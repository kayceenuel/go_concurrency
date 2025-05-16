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

	// Test overwriting existing value
	existed = cache.Put("one", 100)
	if !existed {
		t.Errorf("Key should exist on second Put")
	}

	val, found = cache.Get("one")
	if !found {
		t.Error("Key should exist one second Put")
	}
	if *val != 100 {
		t.Errorf("Expected updated value 100, got %d", *val)
	}

	// Test LRU eviction
	cache.Put("two", 2)
	cache.Put("three", 3)
	cache.Put("four", 4) // Should evict "one" from cache

	_, found = cache.Get("one")
	if found {
		t.Error("Key 'one' should have been evicted")
	}

	// Test that most recenlty items are kept
	cache.Get("two")     // Access "two" to make it most recently used
	cache.Put("five", 5) // should evict "three"

	_, found = cache.Get("three")
	if found {
		t.Error("Key 'three' should have been evicted")
	}

	_, found = cache.Get("two")
	if !found {
		t.Error("Key 'two' should not have been evicted")
	}
}
