package cache

import (
	"sync"
	"zencache/internal/lru"
)

// 支持并发缓存
type cache struct {
	lru      *lru.Cache
	mu       sync.RWMutex
	maxBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.maxBytes, nil)
	}
	c.lru.Add(key, value)
}
func (c *cache) get(key string) (ByteView, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return ByteView{}, false
	}
	value, ok := c.lru.Get(key)
	if !ok {
		return ByteView{}, false
	}
	return value.(ByteView), true
}
