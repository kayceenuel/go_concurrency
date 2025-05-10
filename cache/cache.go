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
	stats      *Statistics          // stats for cache
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

// Get returns the value assiocated with the passed key, and a boolean to indicate
// whether a value was known or not

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.IncrementReads()

	entry, exists := c.items[key]
	if !exists {
		c.stats.IncrementMisses()
		return nil, false
	}

	// Update entry metadata
	entry.accessCount++
	entry.readAfterWrite = true

	// Move to front of LRU list (most recently used)
	c.lruList.moveToFront(key)

	c.stats.IncrementHits()
	return &entry.value, true
}

// GetStatistics returns consistent statistics about the cache
func (c *Cache[K, V]) GetStatistics() Statistics {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a copy of the current statistics
	stats := *c.stats

	// Calculate average access count for current items
	if len(c.items) > 0 {
		totalAccesses := 0
		for _, e := range c.items {
			totalAccesses += e.accessCount
		}
		stats.AverageAccessCount = float64(totalAccesses) / float64(len(c.items))
	}

	currentNeverRead := 0
	for _, e := range c.items {
		if !e.readAfterWrite {
			currentNeverRead++
		}
	}
	stats.CurrentNeverRead = currentNeverRead

	return stats
}

// Node for out doubly linked list
type node[K comparable] struct {
	key  K
	prev *node[K]
	next *node[K]
}

// doublyLinkedlist implement a doubly Linked list for LRU tracking
type doublyLinkedList[K comparable] struct {
	head    *node[K]
	tail    *node[K]
	nodeMap map[K]*node[K] // for O(1) lookups
}

func newDoublyLinkedList[k comparable]() *doublyLinkedList[k] {
	return &doublyLinkedList[k]{
		nodeMap: make(map[k]*node[k]),
	}
}

// addToFront adds a key to the front of the list (most recenly used)
func (l *doublyLinkedList[K]) addToFront(key K) {
	newNode := &node[K]{key: key}

	// Map the key to its node for quick lookup
	l.nodeMap[key] = newNode

	if l.head == nil {
		// List is empty
		l.head = newNode
		l.tail = newNode
		return
	}

	// Add to head
	newNode.next = l.head
	l.head.prev = newNode
	l.head = newNode
}

// moveToFront moves a node to the front of the list
func (l *doublyLinkedList[K]) moveToFront(key K) {
	node, exists := l.nodeMap[key]
	if !exists {
		return
	}

	// Already at front
	if node == l.head {
		return
	}

	// Remove from current position
	if node == l.tail {
		l.tail = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	// Add to front
	node.prev = nil
	node.next = l.head
	l.head.prev = node
	l.head = node
}
