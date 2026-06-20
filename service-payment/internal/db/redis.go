// ─── Redis Storage ─────────────────────────────────────────────
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/model"
)

const (
	paymentKeyPrefix = "payment:"
	defaultTTL       = 72 * time.Hour // 3 days retention
)

// RedisClient wraps a Redis client for payment storage operations.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client from a URL string.
func NewRedisClient(redisURL string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// SavePayment stores a payment as a JSON-encoded hash in Redis.
func (r *RedisClient) SavePayment(ctx context.Context, p *model.Payment) error {
	key := paymentKeyPrefix + p.ID

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal payment: %w", err)
	}

	if err := r.client.Set(ctx, key, data, defaultTTL).Err(); err != nil {
		return fmt.Errorf("save payment to redis: %w", err)
	}

	return nil
}

// GetPayment retrieves a payment by its ID.
func (r *RedisClient) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	key := paymentKeyPrefix + id

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("get payment from redis: %w", err)
	}

	var p model.Payment
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("unmarshal payment: %w", err)
	}

	return &p, nil
}

// UpdatePaymentStatus updates the status and timestamp of an existing payment.
func (r *RedisClient) UpdatePaymentStatus(ctx context.Context, id string, status model.PaymentStatus) error {
	p, err := r.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("payment %s not found", id)
	}

	p.Status = status
	p.UpdatedAt = time.Now().UTC()

	return r.SavePayment(ctx, p)
}

// Close terminates the Redis connection.
func (r *RedisClient) Close() error {
	return r.client.Close()
}
