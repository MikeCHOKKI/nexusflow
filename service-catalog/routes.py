"""
NexusFlow — Catalog Service REST Routes
Endpoints REST utilisés en interne par l'API Gateway.
"""

from __future__ import annotations

import logging
from typing import Optional

from fastapi import APIRouter, HTTPException, Query
from fastapi.responses import Response

from cache import cache
from db import db
from models import Product, ProductCreate, ProductListResponse, ProductUpdate

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/products", tags=["products"])


@router.post("", response_model=Product, status_code=201)
async def create_product(body: ProductCreate) -> Product:
    """Crée un nouveau produit."""
    data = body.model_dump()
    product = await db.create_product(data)
    await cache.invalidate_product(product["id"])
    logger.info("Product created: %s (%s)", product["id"], product["name"])
    return Product(**product)


@router.get("", response_model=ProductListResponse)
async def list_products(
    page: int = Query(default=1, ge=1),
    limit: int = Query(default=20, ge=1, le=100),
    category: Optional[str] = Query(default=None),
    search: Optional[str] = Query(default=None),
) -> ProductListResponse:
    """Liste les produits avec pagination, filtre par catégorie et recherche textuelle."""
    products, total = await db.list_products(
        page=page,
        limit=limit,
        category=category,
        search=search,
    )
    return ProductListResponse(
        products=[Product(**p) for p in products],
        total=total,
        page=page,
        limit=limit,
    )


@router.get("/{product_id}", response_model=Product)
async def get_product(product_id: str) -> Product:
    """Récupère un produit par son ID (cache Redis)."""
    # Cache first
    cached = await cache.get_product(product_id)
    if cached:
        return Product(**cached)

    product = await db.get_product(product_id)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")

    # Cache for subsequent reads
    await cache.set_product(product_id, product)
    return Product(**product)


@router.put("/{product_id}", response_model=Product)
async def update_product(product_id: str, body: ProductUpdate) -> Product:
    """Met à jour partiellement un produit."""
    data = body.model_dump(exclude_unset=True)
    if not data:
        raise HTTPException(status_code=400, detail="No fields to update")

    product = await db.update_product(product_id, data)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")

    await cache.invalidate_product(product_id)
    logger.info("Product updated: %s", product_id)
    return Product(**product)


@router.delete("/{product_id}", status_code=204)
async def delete_product(product_id: str) -> Response:
    """Supprime un produit."""
    deleted = await db.delete_product(product_id)
    if not deleted:
        raise HTTPException(status_code=404, detail="Product not found")

    await cache.invalidate_product(product_id)
    logger.info("Product deleted: %s", product_id)
    return Response(status_code=204)
