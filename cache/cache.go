package cache

import (
	"container/list"
	"sync"
)

// Cache represent a thread safe LRU cache
type Cache struct {
	mu       sync.Mutex               //protects the cache form concurrency access
	capacity int                      //Maximum capacity of the cache
	cache    map[string]*list.Element // Maps key to elements
	list     *list.List               // Doubly linked list LRU order

}

// Pair stores key-value pairs
type pair struct {
	key   int
	value int
}
