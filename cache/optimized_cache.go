package cache

import (
	"sync"
)

// RWMutexCache implements an LRU cache using RWMutex for improved read performance
type RWMutexCache[K comparable, V any] struct {
	mu        sync.RWMutex // Use RWMutex instead of Mutex
	entryLimt int
	items     map[K]*entry[V]
	lrulist   *doublyLinkedList[K]
	stats     *Statistics
}

// NewRWMutexCache creates a new LRU cache with RWMutex
func NewRWMutexCache[K comparable, V any](entryLimit int) *RWMutexCache[K, V] {
	return &RWMutexCache[K, V]{
		entryLimt: entryLimit,
		items:     make(map[K]*entry[V]),
		lrulist:   newDoublyLinkedList[K](),
		stats:     newStatistics(),
	}
}

// Put adds a value to the cache
func (c *RWMutexCache[K, V]) Put(key K, value V) bool {
	c.mu.Lock() // Need exclusive lock for writes
	defer c.mu.Unlock()

	c.stats.IncrementWrites()

	//Implementation same as regular Cache.Put
	existingEntry, exists := c.items[key]
	if exists {
		existingEntry.value = value
		existingEntry.readAfterWrite = false
		c.lrulist.moveToFront(key)
		return true
	}

	if len(c.items) >= c.entryLimt && len(c.items) > 0 {
		lruKey := c.lrulist.removeLast()
		if !c.items[lruKey].readAfterWrite {
			c.stats.IncrementNeverRead()
		}
		delete(c.items, lruKey)
		c.stats.IncrementEvictions()
	}

	c.items[key] = &entry[V]{
		value:          value,
		accessCount:    0,
		readAfterWrite: false,
	}
	c.lrulist.addToFront(key)

	return false
}
