package internal

import (
	"sync"
	"time"
)

type Cache struct {
	mu         sync.Mutex
	interval   time.Duration
	cache_item map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache_item[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ce, ok := c.cache_item[key]
	return ce.val, ok
}

func (c *Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, ci := range c.cache_item {
		entry_time := ci.createdAt
		if time.Since(entry_time) > c.interval {
			delete(c.cache_item, key)
		}
	}
}

func (c *Cache) reaptick(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		c.reapLoop()
	}
}

func NewCache(interval time.Duration) *Cache {
	var new_cache Cache
	new_cache.cache_item = make(map[string]cacheEntry)
	new_cache.interval = interval
	go new_cache.reaptick(new_cache.interval)
	return &new_cache
}
