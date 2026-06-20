# ─── API Gateway ─────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
WORKDIR /app/api-gateway
RUN go mod download

# Proto stubs (shared)
WORKDIR /app
COPY protos/gen/go/ ./protos/gen/go/

# Build service
COPY api-gateway/ ./api-gateway/
WORKDIR /app/api-gateway
RUN CGO_ENABLED=0 go build -o /gateway ./cmd/server

# Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /gateway /gateway
EXPOSE 8080
CMD ["/gateway"]
