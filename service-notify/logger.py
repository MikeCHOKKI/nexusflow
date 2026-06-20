"""Configuration de logging structuré pour le service de notifications.

Format : timestamp, niveau, service name, message, et champs contextuels.
"""

import logging
import sys
from config import settings


def setup_logger(name: str = "nexusflow.notify") -> logging.Logger:
    """Configure et retourne un logger structuré."""
    logger = logging.getLogger(name)
    logger.setLevel(settings.log_level.upper())

    if not logger.handlers:
        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(settings.log_level.upper())

        fmt = (
            "[%(asctime)s] %(levelname)-8s %(name)s "
            "| %(message)s"
        )
        formatter = logging.Formatter(fmt, datefmt="%Y-%m-%d %H:%M:%S")
        handler.setFormatter(formatter)

        logger.addHandler(handler)

    return logger


logger = setup_logger()
