package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheWarehouse map[string]cacheEntry
	mu             sync.Mutex
	interval       time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {

	var ncMap = make(map[string]cacheEntry)

	ncCache := Cache{
		cacheWarehouse: ncMap,
		mu:             sync.Mutex{}, // Not mentioning the Mutex value will behave same,
		// as it will initiliaze with zero value which ready to be used.
		interval: interval,
	}

	go ncCache.reapLoop()

	return ncCache
}

func (c *Cache) Add(key string, val []byte) {
	addCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.mu.Lock()
	c.cacheWarehouse[key] = addCacheEntry
	c.mu.Unlock()

}

func (c *Cache) reapLoop() {
	timeTicker := time.NewTicker(c.interval)
	for range timeTicker.C {
		c.mu.Lock()

		for key, entry := range c.cacheWarehouse {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cacheWarehouse, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {

	c.mu.Lock()
	cacheVal, found := c.cacheWarehouse[key]
	c.mu.Unlock()

	if !found {
		return nil, false
	}

	return cacheVal.val, true
}
