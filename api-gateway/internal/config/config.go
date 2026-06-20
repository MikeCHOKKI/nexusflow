// ─── Configuration ─────────────────────────────────────────────
package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the API Gateway.
type Config struct {
	Port              int
	UserServiceAddr   string
	CatalogServiceAddr string
	OrderServiceAddr  string
	PaymentServiceAddr string
	JWTSecret         string
	RedisURL          string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		Port:              getEnvInt("GATEWAY_PORT", 8080),
		UserServiceAddr:   getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		CatalogServiceAddr: getEnv("CATALOG_SERVICE_ADDR", "localhost:50052"),
		OrderServiceAddr:  getEnv("ORDER_SERVICE_ADDR", "localhost:50053"),
		PaymentServiceAddr: getEnv("PAYMENT_SERVICE_ADDR", "localhost:50054"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
