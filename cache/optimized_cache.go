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
