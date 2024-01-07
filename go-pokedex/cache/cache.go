package cache

import (
	"fmt"
	"sync"
	"time"
)

// Cache entry with creation time (for cleanup) and raw fetched data
type cacheEntry struct {
	createdAt time.Time
	raw       []byte
}

// Cache api calls and specify cleanup interval
type Cache struct {
	Entries  map[string]cacheEntry
	interval time.Duration
	mutex    sync.Mutex
}

// Create a new cache with a cleanup interval
func NewCache(interval time.Duration) Cache {
	return Cache{
		Entries:  map[string]cacheEntry{},
		interval: interval,
	}
}

// Add a new kv pair to our Cache
func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		raw:       value,
	}
}

// Retreive a cache entry (and returns if found)
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	e, found := c.Entries[key]
	return e.raw, found
}

// Cleanup loop that checks age of cache entries and removes
// entries older than Cache.interval
//
// Run this function in a goroutine!
func (c *Cache) CleanupLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.mutex.Lock()

		cleanKeys := []string{}

		for k, v := range c.Entries {
			if t.Sub(v.createdAt) >= c.interval {
				cleanKeys = append(cleanKeys, k)
			}
		}

		for _, k := range cleanKeys {
			fmt.Printf("\n\tCleaning up: %v\n", k)
			delete(c.Entries, k)
		}

		c.mutex.Unlock()
	}

	fmt.Println("Exiting Cleanup")
}
