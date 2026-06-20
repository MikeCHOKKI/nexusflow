// ─── Configuration ─────────────────────────────────────────────
package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the Payment service.
type Config struct {
	GRPCPort    int
	RedisURL    string
	DatabaseURL string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		GRPCPort:    getEnvInt("GRPC_PORT", 50054),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

// getEnv returns the value of the environment variable or the fallback.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt returns the integer value of the environment variable or the fallback.
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
