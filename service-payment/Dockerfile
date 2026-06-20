# ─── Payment Service ─────────────────────────────────────────
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Proto stubs (needed before go mod download because of replace directive)
COPY protos/gen/go/ ./protos/gen/go/

# Cache dependencies
COPY service-payment/go.mod service-payment/go.sum ./service-payment/
WORKDIR /app/service-payment
RUN go mod download

# Build service
COPY service-payment/ ./
RUN CGO_ENABLED=0 go build -o /service ./cmd/server

# Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /service /service
EXPOSE 50054
CMD ["/service"]
