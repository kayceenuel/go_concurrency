package cache

import (
	"sync"
)

// Cache represents a thread safe LRU cache
type Cache[k comparable, V any] struct {
	mu         sync.Mutex           // Mutex for thread safety
	entryLimit int                  // Max num of entries in cache
	items      map[k]*entry[V]      // map of cache entries
	lruList    *doublyLinkedList[K] // list of entries in order of access
	stats      *statistics          // stats for cache
}

// entry represents a cache entry with its value and metadata
type entry[V any] struct {
	value          V    // value of the entry
	accessCount    int  // number of times the entry has been accessed
	readAfterWrite bool // true if write happened after read
}

// NewCahce creates a new LRU cache with the given entry limit
func NewCahce[k comparable, V any](entryLimit int) *Cache[k, V] {
	return &Cache[k, V]{
		entryLimit: entryLimit,
		items:      make(map[k]*entry[V]),
		lruList:    newDoublyLinkedList[k](),
		stats:      newStatistics(),
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether
// a value already existed in the cache fo that key
func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.IncrementWrites()

	// check if the key already exists
	existingEntry, exists := c.items[key]
	if exists {
		// Updata the exisiting key
		existingEntry.value = value
		existingEntry.readAfterWrite = false
		// move to front of LRU list (most recently used)
		c.lruList.moveTofront(key)
		return true
	}

	// key doesn't exist, add it to the cache
	// if we're at capacity, remove the least recently used item

	if len(c.items) >= c.entryLimit && len(c.items) > 0 {
		lruKey := c.lruList.removeLast() // Remove from linked list

		// Update stats before removing
		if !c.items[lruKey].readAfterWrite {
			c.stats.IncrementNeverRead()
		}

		delete(c.items, lruKey) // Remove from map
		c.stats.IncrementEvictions()
	}

	// Add the new entry
	c.items[key] = &entry[V]{
		value:          value,
		accessCount:    0,
		readAfterWrite: false,
	}
	c.lruList.addToFront(key)

	return false
}
