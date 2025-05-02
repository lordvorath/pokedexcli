package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]CacheEntry
	mutex   *sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries: map[string]CacheEntry{},
		mutex:   &sync.Mutex{},
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for t := range ticker.C {
			c.reapLoop(t, interval)
		}
	}()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	// fmt.Printf("CACHE: adding %v\n", key)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, exists bool) {
	/* fmt.Printf("CACHE: getting %v\n", key)
	for key, _ := range c.entries {
		fmt.Printf("-- %v\n", key)
	} */
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) reapLoop(time time.Time, interval time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, val := range c.entries {
		if time.Sub(val.createdAt) > interval {
			delete(c.entries, key)
		}
	}
}
