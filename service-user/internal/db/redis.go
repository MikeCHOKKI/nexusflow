package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient encapsule le client Redis pour cache et blacklist.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initialise la connexion Redis.
func NewRedisClient(redisURL string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &RedisClient{client: client}, nil
}

// Close ferme la connexion Redis.
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// SetToken ajoute un token à la blacklist (durée de vie: expiration du token).
func (r *RedisClient) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, true, ttl).Err()
}

// IsTokenBlacklisted vérifie si un token est dans la blacklist.
func (r *RedisClient) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// CacheGet récupère une valeur depuis le cache.
func (r *RedisClient) CacheGet(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, "cache:"+key).Result()
}

// CacheSet stocke une valeur dans le cache avec une durée de vie.
func (r *RedisClient) CacheSet(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, "cache:"+key, value, ttl).Err()
}

// CacheDel supprime une clé du cache.
func (r *RedisClient) CacheDel(ctx context.Context, key string) error {
	return r.client.Del(ctx, "cache:"+key).Err()
}
