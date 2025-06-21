package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements a Redis-backed cache.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get retrieves a value from Redis by key.
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// Set stores a value in Redis with the given key and TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Close closes the Redis connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
