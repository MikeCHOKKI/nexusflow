"""
NexusFlow — Catalog Service Database
Couche d'accès aux données asyncpg pour les produits.
"""

from __future__ import annotations

import logging
import uuid
from datetime import datetime
from decimal import Decimal
from typing import Optional

import asyncpg

from config import settings

logger = logging.getLogger(__name__)


def _row_to_dict(row: asyncpg.Record) -> dict:
    """Convertit un enregistrement PostgreSQL en dict Python standard.

    Gère les types spécifiques : Decimal → float, UUID → str, datetime → str.
    """
    d = dict(row)
    for key, value in d.items():
        if isinstance(value, Decimal):
            d[key] = float(value)
        elif isinstance(value, datetime):
            d[key] = value.isoformat()
        elif isinstance(value, uuid.UUID):
            d[key] = str(value)
    return d


class Database:
    """Gestionnaire de pool de connexions PostgreSQL."""

    def __init__(self) -> None:
        self.pool: Optional[asyncpg.Pool] = None

    async def connect(self) -> None:
        """Crée le pool de connexions."""
        logger.info("Connecting to database...")
        self.pool = await asyncpg.create_pool(
            settings.database_url,
            min_size=settings.db_pool_min_size,
            max_size=settings.db_pool_max_size,
        )
        logger.info(
            "Database pool created (min=%d, max=%d).",
            settings.db_pool_min_size,
            settings.db_pool_max_size,
        )

    async def close(self) -> None:
        """Ferme le pool de connexions."""
        if self.pool:
            await self.pool.close()
            logger.info("Database pool closed.")

    # ── CRUD ────────────────────────────────────────────────────

    async def create_product(self, data: dict) -> dict:
        """Crée un produit et retourne l'enregistrement complet."""
        async with self.pool.acquire() as conn:  # type: ignore[union-attr]
            row = await conn.fetchrow(
                """
                INSERT INTO catalog.products
                    (name, description, price, currency, category, stock, image_url)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id, name, description, price, currency,
                          category, stock, image_url, is_active,
                          created_at, updated_at
                """,
                data["name"],
                data.get("description"),
                data["price"],
                data.get("currency", "XOF"),
                data["category"],
                data.get("stock", 0),
                data.get("image_url"),
            )
            return _row_to_dict(row)

    async def get_product(self, product_id: str) -> Optional[dict]:
        """Récupère un produit par son ID.

        Retourne None si inexistant.
        """
        async with self.pool.acquire() as conn:  # type: ignore[union-attr]
            row = await conn.fetchrow(
                """
                SELECT id, name, description, price, currency,
                       category, stock, image_url, is_active,
                       created_at, updated_at
                FROM catalog.products
                WHERE id = $1
                """,
                product_id,
            )
            return _row_to_dict(row) if row else None

    async def list_products(
        self,
        page: int = 1,
        limit: int = 20,
        category: Optional[str] = None,
        search: Optional[str] = None,
    ) -> tuple[list[dict], int]:
        """Liste les produits avec pagination, filtre catégorie et recherche.

        Retourne un tuple (liste_produits, total).
        """
        offset = (page - 1) * limit
        conditions: list[str] = []
        params: list = []

        if category:
            conditions.append(f"category = ${len(params) + 1}")
            params.append(category)

        if search:
            conditions.append(f"name ILIKE ${len(params) + 1}")
            params.append(f"%{search}%")

        where_clause = " AND ".join(conditions) if conditions else "TRUE"

        async with self.pool.acquire() as conn:  # type: ignore[union-attr]
            total: int = await conn.fetchval(
                f"SELECT COUNT(*) FROM catalog.products WHERE {where_clause}",
                *params,
            )

            limit_idx = len(params) + 1
            offset_idx = len(params) + 2

            rows = await conn.fetch(
                f"""
                SELECT id, name, description, price, currency,
                       category, stock, image_url, is_active,
                       created_at, updated_at
                FROM catalog.products
                WHERE {where_clause}
                ORDER BY created_at DESC
                LIMIT ${limit_idx} OFFSET ${offset_idx}
                """,
                *params,
                limit,
                offset,
            )

            return [_row_to_dict(r) for r in rows], total or 0

    async def update_product(self, product_id: str, data: dict) -> Optional[dict]:
        """Met à jour partiellement un produit.

        Seuls les champs présents dans ``data`` sont modifiés.
        Retourne le produit mis à jour ou None si inexistant.
        """
        if not data:
            return await self.get_product(product_id)

        valid_fields = {
            "name", "description", "price", "currency",
            "category", "stock", "image_url", "is_active",
        }

        set_clauses: list[str] = []
        params: list = []

        for field, value in data.items():
            if field in valid_fields:
                set_clauses.append(f"{field} = ${len(params) + 1}")
                params.append(value)

        if not set_clauses:
            return await self.get_product(product_id)

        # updated_at mis à jour côté serveur
        set_clauses.append("updated_at = NOW()")

        params.append(product_id)

        async with self.pool.acquire() as conn:  # type: ignore[union-attr]
            row = await conn.fetchrow(
                f"""
                UPDATE catalog.products
                SET {', '.join(set_clauses)}
                WHERE id = ${len(params)}
                RETURNING id, name, description, price, currency,
                          category, stock, image_url, is_active,
                          created_at, updated_at
                """,
                *params,
            )
            return _row_to_dict(row) if row else None

    async def delete_product(self, product_id: str) -> bool:
        """Supprime un produit.

        Retourne True si un enregistrement a été supprimé.
        """
        async with self.pool.acquire() as conn:  # type: ignore[union-attr]
            row = await conn.fetchrow(
                "DELETE FROM catalog.products WHERE id = $1 RETURNING id",
                product_id,
            )
            return row is not None


# Singleton
db = Database()
