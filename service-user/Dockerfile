# ─── User Service ────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY service-user/go.mod service-user/go.sum ./service-user/
WORKDIR /app/service-user
RUN go mod download

# Proto stubs
WORKDIR /app
COPY protos/gen/go/ ./protos/gen/go/

# Build service
COPY service-user/ ./service-user/
WORKDIR /app/service-user
RUN CGO_ENABLED=0 go build -o /service ./cmd/server

# Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /service /service
EXPOSE 50051
CMD ["/service"]
