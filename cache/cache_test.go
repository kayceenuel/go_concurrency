package cache

import (
	"testing"
)

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
