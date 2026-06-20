# ─── Catalog Service ─────────────────────────────────────────
FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends protobuf-compiler && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY service-catalog/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt grpcio-tools

# Generate Python proto stubs
COPY protos/ ./protos/
RUN mkdir -p gen && \
    python -m grpc_tools.protoc --python_out=gen --grpc_python_out=gen -I protos protos/*.proto && \
    touch gen/__init__.py

COPY service-catalog/ .
EXPOSE 8082 50052
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8082"]
