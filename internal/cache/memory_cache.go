package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached value with expiration
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// MemoryCache provides in-memory caching with TTL support
type MemoryCache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewMemoryCache creates a new memory cache with default TTL
func NewMemoryCache(defaultTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
		ttl:   defaultTTL,
	}

	// Start cleanup goroutine
	go cache.startCleanupRoutine()

	return cache
}

// Set stores a value in cache with default TTL
func (mc *MemoryCache) Set(key string, value interface{}) {
	mc.SetWithTTL(key, value, mc.ttl)
}

// SetWithTTL stores a value in cache with custom TTL
func (mc *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from cache
func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	item, exists := mc.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if time.Now().After(item.ExpiresAt) {
		// Don't delete here to avoid deadlock, cleanup routine will handle it
		return nil, false
	}

	return item.Value, true
}

// Delete removes a key from cache
func (mc *MemoryCache) Delete(key string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	delete(mc.items, key)
}

// Clear removes all items from cache
func (mc *MemoryCache) Clear() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.items = make(map[string]*CacheItem)
}

// Size returns the number of items in cache
func (mc *MemoryCache) Size() int {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	return len(mc.items)
}

// startCleanupRoutine periodically removes expired items
func (mc *MemoryCache) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		mc.cleanup()
	}
}

// cleanup removes expired items from cache
func (mc *MemoryCache) cleanup() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	now := time.Now()
	for key, item := range mc.items {
		if now.After(item.ExpiresAt) {
			delete(mc.items, key)
		}
	}
}