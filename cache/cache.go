package cache

import (
	"sync"
)

// Cache represents a thread safe LRU cache
type Cache[K comparable, V any] struct {
	mu         sync.Mutex
	entryLimit int
	items      map[K]*entry[V]
	lruList    *doublyLinkedList[K]
	stats      *statistics
}

// entry represents a cache entry with its value and metadata
type entry[V any] struct {
	value          V    // value of the entry
	accessCount    int  // number of times the entry has been accessed
	readAfterWrite bool // true if write happened after read
}

// NewCahce creates a new LRU cache with the given entry limit
func NewCache[K comparable, V any](entryLimit int) *Cache[K, V] {
	return &Cache[K, V]{
		entryLimit: entryLimit,
		items:      make(map[K]*entry[V]),
		lruList:    newDoublyLinkedList[K](),
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
		c.lruList.moveToFront(key)
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
func (c *Cache[K, V]) GetStatistics() statistics {
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

	// Count never-read entries in current cache
	currentNeverRead := 0
	for _, e := range c.items {
		if !e.readAfterWrite {
			currentNeverRead++
		}
	}
	stats.CurrentNeverRead = currentNeverRead

	return stats
}

// Node for our doubly linked list
type node[K comparable] struct {
	key  K
	prev *node[K]
	next *node[K]
}

// doublyLinkedList implements a doubly linked list for LRU tracking
type doublyLinkedList[K comparable] struct {
	head    *node[K]
	tail    *node[K]
	nodeMap map[K]*node[K] // For O(1) lookups
}

func newDoublyLinkedList[K comparable]() *doublyLinkedList[K] {
	return &doublyLinkedList[K]{
		nodeMap: make(map[K]*node[K]),
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

// removeLast removes and returns the key at the end of the list (least recently used)
func (l *doublyLinkedList[K]) removeLast() K {
	if l.tail == nil {
		var zero K
		return zero
	}

	lastKey := l.tail.key

	//udpate tail
	l.tail = l.tail.prev

	// if this was the last node
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}

	// remove from map
	delete(l.nodeMap, lastKey)

	return lastKey
}
