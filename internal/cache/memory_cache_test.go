package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	// Test setting and getting a value
	cache.Set("test-key", "test-value")
	
	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value, but didn't")
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got %v", value)
	}
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	value, found := cache.Get("non-existent-key")
	if found {
		t.Error("Expected not to find value, but did")
	}

	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	// Set value with very short TTL
	cache.SetWithTTL("expiring-key", "expiring-value", 50*time.Millisecond)

	// Should be available immediately
	value, found := cache.Get("expiring-key")
	if !found || value != "expiring-value" {
		t.Error("Value should be available immediately after setting")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired now
	value, found = cache.Get("expiring-key")
	if found {
		t.Error("Value should be expired")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	cache.Set("delete-key", "delete-value")
	cache.Delete("delete-key")

	value, found := cache.Get("delete-key")
	if found {
		t.Error("Value should be deleted")
	}

	if value != nil {
		t.Errorf("Expected nil value after deletion, got %v", value)
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2, got %d", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", cache.Size())
	}

	// Verify items are actually gone
	_, found := cache.Get("key1")
	if found {
		t.Error("key1 should be cleared")
	}

	_, found = cache.Get("key2")
	if found {
		t.Error("key2 should be cleared")
	}
}

func TestMemoryCache_Size(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	if cache.Size() != 0 {
		t.Error("New cache should be empty")
	}

	cache.Set("key1", "value1")
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}

	cache.Delete("key1")
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after deletion, got %d", cache.Size())
	}
}