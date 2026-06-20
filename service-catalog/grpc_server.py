"""
NexusFlow — Catalog Service gRPC Server
Implémente le service CatalogService défini dans catalog.proto.
"""

from __future__ import annotations

import logging
import sys
from pathlib import Path

import grpc
from grpc import aio as grpc_aio

# ── Import des modules proto générés ────────────────────────────
# Cherche le répertoire 'gen/' produit par protoc.
_PROTO_CANDIDATES = [
    Path(__file__).resolve().parent.parent / "gen",  # ../gen/
    Path("/app/gen"),                                   # Docker runtime
]
for _p in _PROTO_CANDIDATES:
    if _p.exists():
        sys.path.insert(0, str(_p))
        break

import catalog_pb2
import catalog_pb2_grpc
import common_pb2

from cache import cache
from config import settings
from db import db

logger = logging.getLogger(__name__)


# ── Helpers de conversion ───────────────────────────────────────

def _product_to_proto(product: dict) -> catalog_pb2.Product:
    """Convertit un dict produit en message protobuf Product."""
    return catalog_pb2.Product(
        id=product["id"],
        name=product["name"],
        description=product.get("description") or "",
        price=float(product["price"]),
        currency=product.get("currency", "XOF"),
        category=product["category"],
        stock=product.get("stock", 0),
        image_url=product.get("image_url") or "",
        is_active=product.get("is_active", True),
        metadata=common_pb2.Metadata(
            created_at=str(product.get("created_at", "")),
            updated_at=str(product.get("updated_at", "")),
        ),
    )


def _create_request_to_dict(request: catalog_pb2.CreateProductRequest) -> dict:
    return {
        "name": request.name,
        "description": request.description or None,
        "price": request.price,
        "currency": request.currency or "XOF",
        "category": request.category,
        "stock": request.stock,
        "image_url": request.image_url or None,
    }


def _update_request_to_dict(request: catalog_pb2.UpdateProductRequest) -> dict:
    data: dict = {}
    if request.HasField("name"):
        data["name"] = request.name
    if request.HasField("description"):
        data["description"] = request.description or None
    if request.HasField("price"):
        data["price"] = request.price
    if request.HasField("currency"):
        data["currency"] = request.currency
    if request.HasField("category"):
        data["category"] = request.category
    if request.HasField("stock"):
        data["stock"] = request.stock
    if request.HasField("image_url"):
        data["image_url"] = request.image_url or None
    if request.HasField("is_active"):
        data["is_active"] = request.is_active
    return data


# ── Servicer ────────────────────────────────────────────────────

class CatalogServiceServicer(catalog_pb2_grpc.CatalogServiceServicer):
    """Implémentation asynchrone du CatalogService gRPC."""

    async def CreateProduct(
        self,
        request: catalog_pb2.CreateProductRequest,
        context: grpc_aio.ServicerContext,
    ) -> catalog_pb2.Product:
        try:
            data = _create_request_to_dict(request)
            product = await db.create_product(data)
            await cache.invalidate_product(product["id"])
            return _product_to_proto(product)
        except Exception as exc:
            logger.error("gRPC CreateProduct failed: %s", exc)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return catalog_pb2.Product()

    async def GetProduct(
        self,
        request: catalog_pb2.GetProductRequest,
        context: grpc_aio.ServicerContext,
    ) -> catalog_pb2.Product:
        try:
            # Cache first
            product = await cache.get_product(request.id)
            if not product:
                product = await db.get_product(request.id)

            if not product:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"Product {request.id} not found")
                return catalog_pb2.Product()

            # Mettre en cache si pas déjà présent
            await cache.set_product(request.id, product)
            return _product_to_proto(product)
        except Exception as exc:
            logger.error("gRPC GetProduct failed: %s", exc)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return catalog_pb2.Product()

    async def ListProducts(
        self,
        request: catalog_pb2.ListProductsRequest,
        context: grpc_aio.ServicerContext,
    ) -> catalog_pb2.ListProductsResponse:
        try:
            page = max(request.page, 1)
            limit = max(min(request.limit, 100), 1) if request.limit else 20
            category = request.category or None
            search = request.search or None

            products, total = await db.list_products(
                page=page,
                limit=limit,
                category=category,
                search=search,
            )

            return catalog_pb2.ListProductsResponse(
                products=[_product_to_proto(p) for p in products],
                pagination=common_pb2.Pagination(
                    page=page,
                    limit=limit,
                    total=total,
                ),
            )
        except Exception as exc:
            logger.error("gRPC ListProducts failed: %s", exc)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return catalog_pb2.ListProductsResponse()

    async def UpdateProduct(
        self,
        request: catalog_pb2.UpdateProductRequest,
        context: grpc_aio.ServicerContext,
    ) -> catalog_pb2.Product:
        try:
            data = _update_request_to_dict(request)
            product = await db.update_product(request.id, data)

            if not product:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"Product {request.id} not found")
                return catalog_pb2.Product()

            await cache.invalidate_product(request.id)
            return _product_to_proto(product)
        except Exception as exc:
            logger.error("gRPC UpdateProduct failed: %s", exc)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return catalog_pb2.Product()

    async def DeleteProduct(
        self,
        request: catalog_pb2.DeleteProductRequest,
        context: grpc_aio.ServicerContext,
    ) -> catalog_pb2.Empty:
        try:
            deleted = await db.delete_product(request.id)
            if not deleted:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details(f"Product {request.id} not found")
                return catalog_pb2.Empty()

            await cache.invalidate_product(request.id)
            return catalog_pb2.Empty()
        except Exception as exc:
            logger.error("gRPC DeleteProduct failed: %s", exc)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return catalog_pb2.Empty()


# ── Lancement du serveur ────────────────────────────────────────

async def start_grpc_server() -> grpc_aio.Server:
    """Crée, configure et démarre le serveur gRPC asynchrone."""
    server = grpc_aio.server()
    catalog_pb2_grpc.add_CatalogServiceServicer_to_server(
        CatalogServiceServicer(),
        server,
    )
    listen_addr = f"0.0.0.0:{settings.grpc_port}"
    server.add_insecure_port(listen_addr)
    await server.start()
    logger.info("gRPC server listening on %s", listen_addr)
    return server
