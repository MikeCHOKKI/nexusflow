package config

import "os"

// Config regroupe l'ensemble des paramètres du service.
type Config struct {
	GRPCPort   string
	DatabaseURL string
	JWTSecret  string
	RedisURL   string
}

// Load lit les variables d'environnement et retourne une Config.
func Load() *Config {
	return &Config{
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/nexusflow?sslmode=disable"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		RedisURL:   getEnv("REDIS_URL", "redis://localhost:6379/0"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
