"""Gestionnaires d'événements reçus depuis Redis.

Chaque handler reçoit un dict `data` et déclenche les notifications
appropriées (email, log, alerte).
"""

from logger import logger
from notifier import send_email, format_email, send_notification


async def handle_order_created(data: dict) -> None:
    """Une commande vient d'être créée → email de confirmation au client."""
    customer_email = data.get("customer_email", "")
    customer_name = data.get("customer_name", "Client")

    logger.info(
        "Commande créée",
        extra={
            "order_id": data.get("order_id"),
            "customer": customer_email,
        },
    )

    if not customer_email:
        logger.warning("Aucun email client, notification ignorée")
        return

    # Mapper les clés du PHP Order Service vers le template
    email_data = {
        **data,
        "total": data.get("total_amount", data.get("total", "0.00")),
        "customer_name": customer_name,
        "customer_email": customer_email,
        "items": data.get("items", []),
    }

    body = format_email("order_confirmation", email_data)
    await send_email(
        to=customer_email,
        subject=f"Commande #{data.get('order_id', '?')} confirmée",
        body=body,
    )


async def handle_order_status_changed(data: dict) -> None:
    """Statut de commande modifié → notification au client."""
    customer_email = data.get("customer_email", "")
    new_status = data.get("new_status", "inconnu")

    logger.info(
        "Statut commande modifié",
        extra={
            "order_id": data.get("order_id"),
            "new_status": new_status,
            "customer": customer_email,
        },
    )

    if not customer_email:
        logger.warning("Aucun email client, notification ignorée")
        return

    match new_status:
        case "shipped" | "expédiée" | "envoyee":
            body = format_email("order_shipped", data)
            subject = f"Commande #{data.get('order_id', '?')} expédiée"
        case _:
            subject = f"Commande #{data.get('order_id', '?')} : {new_status}"
            body = (
                f"<p>Bonjour,</p>"
                f"<p>Le statut de votre commande "
                f"<strong>#{data.get('order_id', '?')}</strong> "
                f"est passé à : <strong>{new_status}</strong>.</p>"
            )

    await send_email(to=customer_email, subject=subject, body=body)


async def handle_payment_received(data: dict) -> None:
    """Paiement reçu → email de reçu au client."""
    customer_email = data.get("customer_email", "")
    amount = data.get("amount", "0.00")

    logger.info(
        "Paiement reçu",
        extra={
            "order_id": data.get("order_id"),
            "transaction_id": data.get("transaction_id"),
            "amount": amount,
        },
    )

    if not customer_email:
        logger.warning("Aucun email client, notification ignorée")
        return

    body = format_email("payment_receipt", data)
    await send_email(
        to=customer_email,
        subject=f"Reçu de paiement #{data.get('order_id', '?')} — {amount} €",
        body=body,
    )


async def handle_payment_failed(data: dict) -> None:
    """Paiement échoué → log + alerte (pas d'email client par défaut)."""
    order_id = data.get("order_id", "N/A")
    reason = data.get("reason", "Motif inconnu")

    logger.error(
        "Paiement échoué — ALERTE",
        extra={
            "order_id": order_id,
            "reason": reason,
            "transaction_id": data.get("transaction_id"),
        },
    )

    await send_notification(
        channel="admin",
        title=f"Paiement échoué #{order_id}",
        message=f"La transaction {data.get('transaction_id', '?')} "
        f"a échoué pour la commande {order_id} : {reason}",
    )

    customer_email = data.get("customer_email", "")
    if customer_email:
        body = (
            f"<p>Bonjour,</p>"
            f"<p>Le paiement de votre commande "
            f"<strong>#{order_id}</strong> a échoué.</p>"
            f"<p>Raison : {reason}</p>"
            f"<p>Merci de vérifier vos informations bancaires "
            f"et de réessayer.</p>"
        )
        await send_email(
            to=customer_email,
            subject=f"Paiement échoué — Commande #{order_id}",
            body=body,
        )


HANDLER_MAP = {
    "order.created": handle_order_created,
    "order.status_changed": handle_order_status_changed,
    "payment.received": handle_payment_received,
    "payment.failed": handle_payment_failed,
}


async def dispatch(event_type: str, data: dict) -> None:
    """Route un événement vers son handler.

    Args:
        event_type: Type d'événement (ex: "order.created").
        data: Données associées.
    """
    handler = HANDLER_MAP.get(event_type)
    if handler is None:
        logger.warning(
            "Type d'événement inconnu, ignoré",
            extra={"event_type": event_type},
        )
        return
    await handler(data)
