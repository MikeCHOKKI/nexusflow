#!/usr/bin/env python3
"""Notification Service — NexusFlow.

Point d'entrée asynchrone :
- S'abonne au canal Redis "nexusflow:notifications"
- Boucle d'écoute et dispatch vers les handlers
- Serveur HTTP minimal pour healthcheck (GET /health)
"""

import asyncio
import json
from http.server import HTTPServer, BaseHTTPRequestHandler
from threading import Thread

import redis.asyncio as aioredis

from config import settings
from logger import logger
from handlers import dispatch


class HealthHandler(BaseHTTPRequestHandler):
    """Handler HTTP minimal — répond 200 OK sur /health."""

    def do_GET(self) -> None:
        if self.path == "/health":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(b'{"status":"ok","service":"notification"}')
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args) -> None:  # noqa: A002
        logger.debug("HTTP %s", format % args)


def run_health_server() -> None:
    """Lance le serveur healthcheck dans un thread dédié."""
    server = HTTPServer(
        (settings.health_host, settings.health_port),
        HealthHandler,
    )
    logger.info(
        "Healthcheck HTTP démarré",
        extra={"host": settings.health_host, "port": settings.health_port},
    )
    server.serve_forever()


async def listen() -> None:
    """Boucle principale : écoute Redis et dispatch les événements."""
    logger.info(
        "Connexion à Redis",
        extra={"url": settings.redis_url},
    )

    redis_client = aioredis.from_url(
        settings.redis_url,
        decode_responses=True,
    )

    try:
        await redis_client.ping()
        logger.info("Connecté à Redis")
    except Exception as exc:
        logger.error(
            "Impossible de se connecter à Redis",
            extra={"error": str(exc)},
        )
        raise

    pubsub = redis_client.pubsub()
    await pubsub.subscribe(settings.redis_notifications_channel)

    logger.info(
        "Abonné au canal",
        extra={"channel": settings.redis_notifications_channel},
    )

    try:
        async for message in pubsub.listen():
            if message["type"] != "message":
                continue

            raw = message.get("data", "")
            if not raw:
                continue

            try:
                payload = json.loads(raw)
            except json.JSONDecodeError as exc:
                logger.warning(
                    "Message JSON invalide, ignoré",
                    extra={"error": str(exc), "raw": raw[:200]},
                )
                continue

            event_type = payload.get("type", "unknown")
            data = payload.get("data", {})

            logger.info(
                "Événement reçu",
                extra={"type": event_type, "order_id": data.get("order_id")},
            )

            try:
                await dispatch(event_type, data)
            except Exception as exc:
                logger.error(
                    "Erreur lors du traitement de l'événement",
                    extra={
                        "type": event_type,
                        "error": str(exc),
                    },
                )
    except asyncio.CancelledError:
        logger.info("Boucle d'écoute interrompue")
    finally:
        await pubsub.unsubscribe(settings.redis_notifications_channel)
        await redis_client.close()
        logger.info("Déconnecté de Redis")


async def main() -> None:
    logger.info(
        "Démarrage Notification Service",
        extra={"channel": settings.redis_notifications_channel},
    )

    health_thread = Thread(target=run_health_server, daemon=True)
    health_thread.start()

    await listen()


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Arrêt demandé par l'utilisateur")
    except Exception as exc:
        logger.error(
            "Arrêt inattendu",
            extra={"error": str(exc)},
        )
        raise
