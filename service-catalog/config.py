"""
NexusFlow — Catalog Service Configuration
Charge les variables d'environnement pour le service.
"""

from __future__ import annotations

import os
from dataclasses import dataclass, field


@dataclass(frozen=True)
class Settings:
    """Configuration immutable du service catalog."""

    # Ports
    grpc_port: int = int(os.getenv("GRPC_PORT", "50053"))

    # Bases de données
    database_url: str = os.getenv(
        "DATABASE_URL",
        "postgresql://nexusflow:nexusflow@localhost:5432/nexusflow",
    )
    redis_url: str = os.getenv("REDIS_URL", "redis://localhost:6379/0")

    # Pool DB
    db_pool_min_size: int = int(os.getenv("DB_POOL_MIN_SIZE", "2"))
    db_pool_max_size: int = int(os.getenv("DB_POOL_MAX_SIZE", "10"))

    # Cache
    cache_ttl: int = int(os.getenv("CACHE_TTL", "300"))  # 5 minutes

    # Logging
    log_level: str = os.getenv("LOG_LEVEL", "INFO")


settings = Settings()
