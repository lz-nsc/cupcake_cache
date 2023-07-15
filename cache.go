package cupcake_cache

import (
	"sync"

	"github.com/lz-nsc/cupcake_cache/lru"
)

type cache struct {
	mu    sync.Mutex
	cache *lru.LRUCache
	size  int64 // maximum bytes of the cache
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = lru.New(c.size, nil)
	}
	c.cache.Add(key, value)
}

func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}

	if v, ok := c.cache.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
