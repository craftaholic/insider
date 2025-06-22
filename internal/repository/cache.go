package repository

import (
	"time"

	"github.com/craftaholic/insider/internal/domain/interfaces"
	"github.com/go-redis/redis"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) interfaces.CacheRepository {
	return &CacheRepository{
		client: redisClient,
	}
}

func (cr *CacheRepository) Get(key string) ([]byte, error) {
	return cr.client.Get(key).Bytes()
}

func (cr *CacheRepository) Set(key string, value []byte, ttl time.Duration) error {
	return cr.client.Set(key, value, ttl).Err()
}
