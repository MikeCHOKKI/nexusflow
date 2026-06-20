"""
NexusFlow — Catalog Service
Point d'entrée FastAPI.
- Initialise la base de données (asyncpg) et Redis
- Monte les routes REST et le serveur gRPC en arrière-plan
- Endpoint GET /health
"""

from __future__ import annotations

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from cache import cache
from config import settings
from db import db
from grpc_server import start_grpc_server
from routes import router

# ── Logging structuré ───────────────────────────────────────────
logging.basicConfig(
    level=getattr(logging, settings.log_level.upper(), logging.INFO),
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%dT%H:%M:%S%z",
)
logger = logging.getLogger("catalog")

# Référence globale au serveur gRPC pour l'arrêt
_grpc_server = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Gère le cycle de vie du service : startup / shutdown."""
    global _grpc_server

    # ── Startup ─────────────────────────────────────────────
    logger.info("═══ Catalog Service starting ═══")
    logger.info("gRPC port: %d  |  Log level: %s", settings.grpc_port, settings.log_level)

    await db.connect()
    await cache.connect()

    _grpc_server = await start_grpc_server()

    yield

    # ── Shutdown ────────────────────────────────────────────
    logger.info("═══ Catalog Service shutting down ═══")
    if _grpc_server:
        await _grpc_server.stop(grace=5)
        logger.info("gRPC server stopped.")
    await cache.close()
    await db.close()
    logger.info("Catalog Service stopped.")


# ── Application FastAPI ────────────────────────────────────────
app = FastAPI(
    title="NexusFlow Catalog Service",
    description="Microservice de gestion des produits du catalogue",
    version="1.0.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
)

# CORS — ouvert en interne derrière l'API Gateway
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Montage des routes REST
app.include_router(router, prefix="/api/v1")


@app.get("/health")
async def health() -> dict:
    """Endpoint de health-check."""
    return {
        "status": "ok",
        "service": "catalog",
        "grpc_port": settings.grpc_port,
    }
