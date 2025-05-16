package cache

import (
	"testing"
)

func TestBasicCacheOperations(t *testing.T) {
	cache := NewCache[string, int](3)

	// Test Put and Get
	existed := cache.Put("one", 1)
	if existed {
		t.Error("key shouldn't exist on first put")
	}

	val, found := cache.Get("one")
	if !found {
		t.Error("key should exist after Put")
	}
	if *val != 1 {
		t.Errorf("Expected value 1, got %d", *val)
	}
}
