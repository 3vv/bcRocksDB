package rocksdb

//#include "api.h"
import "C"

// Cache is a cache used to store data read from data in memory.
type Cache struct {
	c *C.rocksdb_cache_t
}

// NewLRUCache creates a new LRU Cache object with the capacity given.
func NewLRUCache(capacity int) *Cache {
	return NewNativeCache(C.rocksdb_cache_create_lru(C.size_t(capacity)))
}

// NewNativeCache creates a Cache object.
func NewNativeCache(c *C.rocksdb_cache_t) *Cache {
	return &Cache{c}
}

// GetUsage returns the Cache memory usage.
func (c *Cache) GetUsage() int {
	return int(C.rocksdb_cache_get_usage(c.c))
}

// GetPinnedUsage returns the Cache pinned memory usage.
func (c *Cache) GetPinnedUsage() int {
	return int(C.rocksdb_cache_get_pinned_usage(c.c))
}

// Destroy deallocates the Cache object.
func (c *Cache) Destroy() {
	C.rocksdb_cache_destroy(c.c)
	c.c = nil
}