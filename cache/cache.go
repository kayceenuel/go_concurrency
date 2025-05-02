package cache

import (
	"container/list"
	"sync"
)

// Cache represent a thread safe LRU cache
type LRUCache struct {
	mu       sync.Mutex               //protects the cache form concurrency access
	capacity int                      //Maximum capacity of the cache
	cache    map[string]*list.Element // Maps key to elements
	list     *list.List               // Doubly linked list LRU order

}

// Pair stores key-value pairs
type Pair struct {
	key   string
	value interface{}
}

// NewCache initializes the LRUCache with a given capacity
func NewCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// "Get" retrieves a value the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*Pair).value, true
	}
	return nil, false

}

func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.value.(*Pair).value = value
		return
	}
	if c.list.Len() >= c.cap {
		//Evict the LRU Items
		lru := c.list.Back()
		if lru != nil {
			c.list.Remove(lru)
			delete(c.cache, lru.Value.(*Pair).key)
		}
	}
	pair := &Pair{key: key, value: value}
	elem := c.list.PushFront(pair)
	c.cache[key] = elem
}
