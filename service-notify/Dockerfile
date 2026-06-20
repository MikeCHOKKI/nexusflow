# ─── Notification Service ────────────────────────────────────
FROM python:3.12-slim
WORKDIR /app

COPY service-notify/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY service-notify/ .
CMD ["python", "main.py"]
