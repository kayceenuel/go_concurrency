package cache

import (
	"sync"
	"testing"
)

// TestCacheBasic checks basic LRU functionality (insertion, retrieval , and eviction).
func TestCacheBasic(t *testing.T) {
	c := NewCache(2)
	c.Put("key1", 1)
	c.Put("key2", 2)
	if val, ok := c.Get("key1"); !ok || val != 1 {
		t.Errorf("Expected key1 to return 1, got %v", val)
	}
	c.Put("key3", 3) // should evict key2
	if _, ok := c.Get("key2"); ok {
		t.Error("Expected key2 to be evicted")
	}
}

// TestCacheConcurrency tests thread safety with multiple goroutines accessing the cache concurrently.
func TestCacheConcurrency(t *testing.T) {
	c := NewCache(10)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('A' + i%26))
			c.Put(key, i)
			if val, ok := c.Get(key); !ok || val.(int) != i {
				t.Errorf("Expected %s to return %d, got %v", key, i, val)
			}
		}(i)
	}
	wg.Wait()
}
