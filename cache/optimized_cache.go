package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
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

func (c *RWMutexCache[K, V]) Get(key K) (*V, bool) {
	// First try a read lock for the lookup
	c.mu.RLock()
	entry, exists := c.items[key]
	c.mu.RUnlock()

	c.stats.IncrementReads()

	if !exists {
		c.stats.IncrementMisses()
		return nil, false
	}

	// Now we need to update the LRU list and Metadata, which requires a write lock
	c.mu.Lock()
	// Double-check the entry still exists (it might have been evicted in between locks)
	entry, stillExist := c.items[key]
	if !stillExist {
		c.mu.Unlock()
		c.stats.IncrementMisses()
		return nil, false
	}

	// Update entry metadata
	entry.accessCount++
	entry.readAfterWrite = true
	c.lrulist.moveToFront(key)
	c.stats.IncrementHits()

	// Make a copy of the value to return
	result := entry.value
	c.mu.Unlock()

	return &result, true
}

// GetStatistics returns consistent statistics about the cache
func (c *RWMutexCache[K, V]) GetStatistics() Statistics {
	c.mu.RLock() // Use read lock for statistics
	defer c.mu.RUnlock()

	// Similar to regular Cache.GetStatistics
	stats := *c.stats

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

// SharededCache implements an LRU cache with mutiple shards for reduced lock contention
type ShardedCache[K comparable, V any] struct {
	shards     []*Cache[K, V]
	shardCount int
	shardMask  int
	stats      atomic.Pointer[Statistics]
}

func NewShardedCache[K comparable, V any](entryLimit int, shardCount int) *ShardedCache[K, V] {
	// Make sure shardCount is a power of 2 for efficient modulo
	if shardCount&(shardCount-1) != 0 {
		// Find next power of 2
		shardCount--
		shardCount |= shardCount >> 1
		shardCount |= shardCount >> 2
		shardCount |= shardCount >> 4
		shardCount |= shardCount >> 8
		shardCount |= shardCount >> 16
		shardCount++
	}

	entriesPerShard := entryLimit / shardCount
	if entriesPerShard < 1 {
		entriesPerShard = 1
	}

	cache := &ShardedCache[K, V]{
		shards:     make([]*Cache[K, V], shardCount),
		shardCount: shardCount,
		shardMask:  shardCount,
	}

	for i := 0; i < shardCount; i++ {
		cache.shards[i] = NewCache[K, V](entriesPerShard)
	}

	// Initialize stats
	initialStats := newStatistics()
	cache.stats.Store((initialStats))

	return cache
}

// getShard returns the apporpriate shard for a key
func (c *ShardedCache[K, V]) getShard(key K) *Cache[K, V] {
	// we need a hash function for the key
	// This is a simple one; a better implementation would use a good hash function
	h := anyToHash(key)
	return c.shards[h&c.shardMask]
}

// anyToHash converts any comparable to a unit64 hash
// This is a very simple hash function
func anyToHash[K comparable](key K) int {
	// Using a simple hash function based on fmt.Sprintf and basic string hashing
	s := fmt.Sprintf("%v", key)
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	if h < 0 {
		h = -h
	}
	return h
}

// Put adds a value to the cache
func (c *ShardedCache[K, V]) Put(key K, value V) bool {
	return c.getShard(key).Put(key, value)
}

// Get retrieves a value from the cache
func (c *ShardedCache[K, V]) Get(key K) (*V, bool) {
	return c.getShard(key).Get(key)
}

// GetStatistics returns aggregate statistics about the cache
func (c *ShardedCache[K, V]) GetStatistics() Statistics {
	// Collect stats from all shards
	aggregateStats := Statistics{}

	for _, shard := range c.shards {
		shardStats := shard.GetStatistics()

		// Aggregate counters
		aggregateStats.Reads += shardStats.Reads
		aggregateStats.Writes += shardStats.Writes
		aggregateStats.Hits += shardStats.Hits
		aggregateStats.Misses += shardStats.Misses
		aggregateStats.Evictions += shardStats.Evictions
		aggregateStats.NeverReadCount += shardStats.NeverReadCount
		aggregateStats.CurrentNeverRead += shardStats.CurrentNeverRead
	}

	// Calculate hit rate
	if aggregateStats.Reads > 0 {
		aggregateStats.AverageAccessCount = float64(aggregateStats.Hits) / float64(c.shardCount)
	}

	return aggregateStats
}
