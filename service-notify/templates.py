"""Templates d'emails HTML pour les notifications NexusFlow.

Chaque template est une fonction qui prend un dict `data` et retourne
un body HTML complet (string).
"""


def order_confirmation(data: dict) -> str:
    order_id = data.get("order_id", "N/A")
    customer_name = data.get("customer_name", "Client")
    total = data.get("total", "0.00")
    items = data.get("items", [])

    items_html = ""
    for item in items:
        qty = item.get('quantity', 1)
        unit_price = item.get('unit_price', 0)
        items_html += (
            f"<tr>"
            f"<td style='padding:8px;border-bottom:1px solid #eee;'>{item.get('product_name', '')}</td>"
            f"<td style='padding:8px;border-bottom:1px solid #eee;text-align:center;'>{qty}</td>"
            f"<td style='padding:8px;border-bottom:1px solid #eee;text-align:right;'>{float(unit_price):,.0f} FCFA</td>"
            f"</tr>"
        )

    return f"""<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#f8f9fa;border-radius:12px;padding:24px;">
    <h2 style="margin-top:0;color:#1a1a2e;">✅ Commande confirmée</h2>
    <p>Bonjour <strong>{customer_name}</strong>,</p>
    <p>Votre commande <strong>#{order_id}</strong> a bien été confirmée.</p>
    <table style="width:100%;border-collapse:collapse;margin:16px 0;">
      <thead>
        <tr style="background:#e9ecef;">
          <th style="padding:8px;text-align:left;">Article</th>
          <th style="padding:8px;text-align:center;">Qté</th>
          <th style="padding:8px;text-align:right;">Prix</th>
        </tr>
      </thead>
      <tbody>
        {items_html if items else "<tr><td colspan='3' style='padding:8px;text-align:center;color:#888;'>Aucun article</td></tr>"}
      </tbody>
    </table>
    <p style="font-size:18px;font-weight:bold;text-align:right;margin:0;">Total : {float(total):,.0f} FCFA</p>
    <hr style="border:none;border-top:1px solid #dee2e6;margin:20px 0;">
    <p style="color:#6c757d;font-size:13px;">NexusFlow — Service client</p>
  </div>
</body>
</html>"""


def order_shipped(data: dict) -> str:
    order_id = data.get("order_id", "N/A")
    customer_name = data.get("customer_name", "Client")
    tracking = data.get("tracking_url", "")
    carrier = data.get("carrier", "Transporteur")

    tracking_block = ""
    if tracking:
        tracking_block = (
            f'<p style="margin:16px 0;">'
            f'<a href="{tracking}" '
            f'style="display:inline-block;background:#1a73e8;color:#fff;'
            f'text-decoration:none;padding:12px 24px;border-radius:6px;'
            f'font-weight:600;">Suivre mon colis</a></p>'
        )

    return f"""<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#f8f9fa;border-radius:12px;padding:24px;">
    <h2 style="margin-top:0;color:#1a1a2e;">📦 Commande expédiée</h2>
    <p>Bonjour <strong>{customer_name}</strong>,</p>
    <p>Votre commande <strong>#{order_id}</strong> a été expédiée via <strong>{carrier}</strong>.</p>
    {tracking_block}
    <hr style="border:none;border-top:1px solid #dee2e6;margin:20px 0;">
    <p style="color:#6c757d;font-size:13px;">NexusFlow — Service expédition</p>
  </div>
</body>
</html>"""


def payment_receipt(data: dict) -> str:
    order_id = data.get("order_id", "N/A")
    customer_name = data.get("customer_name", "Client")
    amount = data.get("amount", "0.00")
    payment_method = data.get("payment_method", "Carte bancaire")
    transaction_id = data.get("transaction_id", "N/A")

    return f"""<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#f8f9fa;border-radius:12px;padding:24px;">
    <h2 style="margin-top:0;color:#1a1a2e;">💰 Paiement reçu</h2>
    <p>Bonjour <strong>{customer_name}</strong>,</p>
    <p>Nous confirmons la réception de votre paiement pour la commande <strong>#{order_id}</strong>.</p>
    <table style="width:100%;border-collapse:collapse;margin:16px 0;">
      <tr><td style="padding:8px;color:#555;">Montant</td><td style="padding:8px;font-weight:bold;">{float(amount):,.0f} FCFA</td></tr>
      <tr><td style="padding:8px;color:#555;">Mode de paiement</td><td style="padding:8px;">{payment_method}</td></tr>
      <tr><td style="padding:8px;color:#555;">Transaction</td><td style="padding:8px;font-family:monospace;font-size:13px;">{transaction_id}</td></tr>
    </table>
    <hr style="border:none;border-top:1px solid #dee2e6;margin:20px 0;">
    <p style="color:#6c757d;font-size:13px;">NexusFlow — Service facturation</p>
  </div>
</body>
</html>"""
