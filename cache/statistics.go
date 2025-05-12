package cache

import (
	"sync/atomic"
)

// Statistics tracks cache usage statistics
type Statistics struct {
	Reads  int64 // Total number of read operations
	Writes int64 // Total number of write operations
	Hits   int64 // Number of cache hits
}

// newStatistics creates a new statistics tracker
func newStatistics() *Statistics {
	return &Statistics{}
}

func (s *Statistics) IncrementReads() {
	atomic.AddInt64(&s.Reads, 1)
}

func (s *Statistics) IncrementWrites() {
	atomic.AddInt64(&s.Writes, 1)
}

func (s *Statistics) IncrementHits() {
	atomic.AddInt64(&s.Hits, 1)
}
