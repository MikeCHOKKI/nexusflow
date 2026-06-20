"""
NexusFlow — Catalog Service Models
Pydantic v2 models pour les produits du catalogue.
"""

from __future__ import annotations

from typing import Optional

from pydantic import BaseModel, Field


class Product(BaseModel):
    """Modèle de réponse complet pour un produit."""

    id: str
    name: str
    description: Optional[str] = None
    price: float
    currency: str = "XOF"
    category: str
    stock: int = 0
    image_url: Optional[str] = None
    is_active: bool = True
    created_at: Optional[str] = None
    updated_at: Optional[str] = None

    model_config = {"from_attributes": True, "populate_by_name": True}


class ProductCreate(BaseModel):
    """Payload de création d'un produit."""

    name: str = Field(..., min_length=1, max_length=255)
    description: Optional[str] = Field(default=None, max_length=2000)
    price: float = Field(..., gt=0)
    currency: str = Field(default="XOF", pattern=r"^(XOF|EUR|USD)$")
    category: str = Field(..., min_length=1, max_length=100)
    stock: int = Field(default=0, ge=0)
    image_url: Optional[str] = Field(default=None, max_length=500)


class ProductUpdate(BaseModel):
    """Payload de mise à jour partielle d'un produit."""

    name: Optional[str] = Field(default=None, min_length=1, max_length=255)
    description: Optional[str] = Field(default=None, max_length=2000)
    price: Optional[float] = Field(default=None, gt=0)
    currency: Optional[str] = Field(default=None, pattern=r"^(XOF|EUR|USD)$")
    category: Optional[str] = Field(default=None, min_length=1, max_length=100)
    stock: Optional[int] = Field(default=None, ge=0)
    image_url: Optional[str] = Field(default=None, max_length=500)
    is_active: Optional[bool] = None


class ProductListResponse(BaseModel):
    """Réponse paginée pour la liste des produits."""

    products: list[Product]
    total: int
    page: int
    limit: int
