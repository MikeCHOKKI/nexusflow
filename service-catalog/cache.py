"""
NexusFlow — Catalog Service Cache
Cache Redis pour les produits les plus consultés (TTL 5 min).
"""

from __future__ import annotations

import json
import logging
from typing import Optional

import redis.asyncio as aioredis

from config import settings

logger = logging.getLogger(__name__)


class ProductCache:
    """Cache Redis pour les produits du catalogue."""

    def __init__(self) -> None:
        self._redis: Optional[aioredis.Redis] = None

    async def connect(self) -> None:
        """Initialise la connexion Redis."""
        logger.info("Connecting to Redis at %s ...", settings.redis_url)
        self._redis = aioredis.from_url(
            settings.redis_url,
            decode_responses=True,
        )
        await self._redis.ping()
        logger.info("Redis connected.")

    async def close(self) -> None:
        """Ferme la connexion Redis."""
        if self._redis:
            await self._redis.close()
            logger.info("Redis connection closed.")

    # ── Clés ────────────────────────────────────────────────────

    @staticmethod
    def _product_key(product_id: str) -> str:
        return f"catalog:product:{product_id}"

    # ── Opérations ──────────────────────────────────────────────

    async def get_product(self, product_id: str) -> Optional[dict]:
        """Récupère un produit depuis le cache.

        Retourne le dict produit ou None si non trouvé / expiré.
        """
        if not self._redis:
            return None
        try:
            data = await self._redis.get(self._product_key(product_id))
            if data:
                return json.loads(data)
            return None
        except Exception:
            logger.warning("Cache get failed for product %s", product_id, exc_info=True)
            return None

    async def set_product(self, product_id: str, data: dict) -> None:
        """Stocke un produit dans le cache avec TTL configuré."""
        if not self._redis:
            return
        try:
            await self._redis.setex(
                self._product_key(product_id),
                settings.cache_ttl,
                json.dumps(data, default=str),
            )
        except Exception:
            logger.warning("Cache set failed for product %s", product_id, exc_info=True)

    async def invalidate_product(self, product_id: str) -> None:
        """Invalide le cache d'un produit spécifique."""
        if not self._redis:
            return
        try:
            await self._redis.delete(self._product_key(product_id))
        except Exception:
            logger.warning("Cache invalidate failed for product %s", product_id, exc_info=True)


# Singleton
cache = ProductCache()
