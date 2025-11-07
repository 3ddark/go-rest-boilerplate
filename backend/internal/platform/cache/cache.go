package cache

import (
	"sync"
	"time"
)

type ICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, duration time.Duration)
	Delete(key string)
}

type InMemoryCache struct {
	items map[string]item
	mu    sync.RWMutex
}

type item struct {
	value      interface{}
	expiration int64
}

func NewInMemoryCache(cleanupInterval time.Duration) ICache {
	c := &InMemoryCache{
		items: make(map[string]item),
	}
	go c.janitor(cleanupInterval)
	return c
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, found := c.items[key]
	if !found || time.Now().UnixNano() > it.expiration {
		return nil, false
	}

	return it.value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: expiration,
	}
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *InMemoryCache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, it := range c.items {
			if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
