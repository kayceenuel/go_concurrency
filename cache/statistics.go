package cache

import (
	"sync/atomic"
)

// Statistics tracks cache usage statistics
type Statistics struct {
	Reads int64 // Total number of read operations
}

// newStatistics creates a new statistics tracker
func newStatistics() *Statistics {
	return &Statistics{}
}

func (s *Statistics) IncrementReads() {
	atomic.AddInt64(&s.Reads, 1)
}
