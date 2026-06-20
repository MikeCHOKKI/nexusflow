"""Envoi d'emails et dispatch de notifications.

Utilise aiosmtplib pour l'envoi SMTP asynchrone (via MailHog en dev)
et logge les notifications pour les autres canaux.
"""

import aiosmtplib
from email.message import EmailMessage

from config import settings
from logger import logger
from templates import order_confirmation, order_shipped, payment_receipt


SMTP_TIMEOUT = 10


async def send_email(to: str, subject: str, body: str) -> bool:
    """Envoie un email HTML via SMTP.

    Args:
        to: Destinataire (email).
        subject: Sujet du message.
        body: Corps HTML.

    Returns:
        True si envoyé, False en cas d'erreur (loggée, pas de crash).
    """
    msg = EmailMessage()
    msg["From"] = settings.smtp_from
    msg["To"] = to
    msg["Subject"] = subject
    msg.set_content(body, subtype="html")

    try:
        await aiosmtplib.send(
            msg,
            hostname=settings.smtp_host,
            port=settings.smtp_port,
            timeout=SMTP_TIMEOUT,
        )
        logger.info(
            "Email envoyé",
            extra={"to": to, "subject": subject},
        )
        return True
    except Exception as exc:
        logger.error(
            "Échec envoi email vers %s (sujet: %s) : %s",
            to,
            subject,
            str(exc),
        )
        return False


TEMPLATE_MAP = {
    "order_confirmation": order_confirmation,
    "order_shipped": order_shipped,
    "payment_receipt": payment_receipt,
}


def format_email(template_name: str, data: dict) -> str:
    """Construit un email HTML à partir d'un template et de données.

    Args:
        template_name: Nom du template (clé dans TEMPLATE_MAP).
        data: Données à injecter dans le template.

    Returns:
        Corps HTML complet.

    Raises:
        ValueError: Si le template est inconnu.
    """
    tpl = TEMPLATE_MAP.get(template_name)
    if tpl is None:
        raise ValueError(f"Template inconnu : {template_name}")
    return tpl(data)


async def send_notification(
    channel: str, title: str, message: str, webhook_url: str | None = None
) -> None:
    """Dispatch d'une notification vers un ou plusieurs canaux.

    Pour l'instant : log structuré uniquement.
    Pourrait être étendu vers webhook, Slack, etc.

    Args:
        channel: Canal cible (ex: "admin", "support").
        title: Titre de la notification.
        message: Corps du message.
        webhook_url: URL optionnelle pour webhook HTTP.
    """
    logger.info(
        "Notification",
        extra={
            "channel": channel,
            "title": title,
            "message": message,
            "webhook": webhook_url or "none",
        },
    )

    if webhook_url:
        try:
            import httpx

            async with httpx.AsyncClient(timeout=5) as client:
                payload = {"title": title, "message": message, "channel": channel}
                resp = await client.post(webhook_url, json=payload)
                resp.raise_for_status()
                logger.info("Webhook envoyé", extra={"url": webhook_url})
        except ImportError:
            logger.warning("httpx non installé, webhook ignoré")
        except Exception as exc:
            logger.error(
                "Échec webhook",
                extra={"url": webhook_url, "error": str(exc)},
            )
