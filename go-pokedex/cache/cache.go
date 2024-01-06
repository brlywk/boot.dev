package cache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	raw       []byte
}

type Cache struct {
	Entries  map[string]cacheEntry
	interval time.Duration
	mutex    sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	return Cache{
		Entries:  map[string]cacheEntry{},
		interval: interval,
	}
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		raw:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	e, found := c.Entries[key]
	return e.raw, found
}

func (c *Cache) CleanupLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Println("\n\tTime to clean up a little bit!")
		// fmt.Printf("\n\tCache:\n%v\n", c)

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
