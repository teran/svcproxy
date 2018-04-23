package redis

import (
	"context"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/acme/autocert"
)

var _ autocert.Cache = &Cache{}

// Cache type for Redis cache implementation
type Cache struct {
	client *redis.Client
}

// NewCache returns new Cache instance
func NewCache(client *redis.Client) (autocert.Cache, error) {
	return &Cache{
		client: client,
	}, nil
}

// Get retrieves data by key from Redis
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(key).Result()
	if err != nil && err == redis.Nil {
		return nil, autocert.ErrCacheMiss
	}

	return []byte(data), err
}

// Put stores data with key to Redis
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	return c.client.Set(key, string(data), 0).Err()
}

// Delete delete data by key from Redis
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(key).Err()
}
