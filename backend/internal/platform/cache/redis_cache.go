package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache is the Redis implementation of the ICache interface.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(client *redis.Client) ICache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Get retrieves an item from the cache.
func (c *RedisCache) Get(key string) (interface{}, bool) {
	val, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false // Key does not exist
	} else if err != nil {
		// Log the error but treat as a cache miss
		return nil, false
	}
	return val, true
}

// Set adds an item to the cache for a specified duration.
func (c *RedisCache) Set(key string, value interface{}, duration time.Duration) {
	// Redis client handles non-byte values, but byte slices are common.
	err := c.client.Set(c.ctx, key, value, duration).Err()
	if err != nil {
		// Log the error
	}
}

// Delete removes an item from the cache.
func (c *RedisCache) Delete(key string) {
	c.client.Del(c.ctx, key)
}
