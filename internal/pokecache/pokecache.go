package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	tt := time.NewTicker(interval)
	for range tt.C {
		cutoff := time.Now().UTC().Add(-interval)
		c.mu.Lock()
		for key, val := range c.entries {
			if !val.createdAt.After(cutoff) {
				fmt.Printf("%v cache deleting: %v\n", time.Now(), key)
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
