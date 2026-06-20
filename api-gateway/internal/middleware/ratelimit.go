// ─── Rate Limiting Middleware (Token Bucket via Redis) ─────────
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRequestsPerMinute = 100
	rateLimitWindow      = time.Minute
)

// RateLimit returns an HTTP middleware that enforces a token-bucket rate limit
// per client IP using Redis. Returns HTTP 429 when the limit is exceeded.
func RateLimit(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			key := fmt.Sprintf("ratelimit:%s", ip)

			count, err := rdb.Incr(r.Context(), key).Result()
			if err != nil {
				// Redis unavailable — allow the request but log.
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				// First request in this window — set expiry.
				rdb.Expire(r.Context(), key, rateLimitWindow)
			}

			if count > maxRequestsPerMinute {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				writeRateLimitError(w, "rate limit exceeded: 100 req/min")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// NewRedisClient creates a new Redis client from a connection URL.
func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis parse url: %w", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return rdb, nil
}

func extractIP(r *http.Request) string {
	// Respect X-Forwarded-For when behind a reverse proxy.
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return real
	}
	return r.RemoteAddr
}

func writeRateLimitError(w http.ResponseWriter, msg string) {
	w.Write([]byte(`{"success":false,"data":null,"error":"` + msg + `"}`))
}
