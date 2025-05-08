package cache

import (
	"sync"
)

// Cache represents a thread safe LRU cache
type Cahce[k comparable, V any] struct {
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
func NewCahce[k comparable, V any](entryLimit int) *Cahce[k, V] {
	return &Cahce[k, V]{
		entryLimit: entryLimit,
		items:      make(map[k]*entry[V]),
		lruList:    newDoublyLinkedList[k](),
		stats:      newStatistics(),
	}
}
