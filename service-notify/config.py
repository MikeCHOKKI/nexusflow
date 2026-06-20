"""Configuration du service de notifications NexusFlow.

Toutes les variables sont chargées depuis l'environnement
avec des valeurs par défaut pour le développement local.
"""

from os import getenv
from dotenv import load_dotenv

load_dotenv()


class Settings:
    redis_url: str = getenv("REDIS_URL", "redis://localhost:6379/0")

    smtp_host: str = getenv("SMTP_HOST", "localhost")
    smtp_port: int = int(getenv("SMTP_PORT", "1025"))
    smtp_from: str = getenv("SMTP_FROM", "notifications@nexusflow.local")

    log_level: str = getenv("LOG_LEVEL", "INFO")

    redis_notifications_channel: str = getenv(
        "REDIS_NOTIFICATIONS_CHANNEL", "nexusflow:notifications"
    )

    health_host: str = getenv("HEALTH_HOST", "0.0.0.0")
    health_port: int = int(getenv("HEALTH_PORT", "50055"))


settings = Settings()
