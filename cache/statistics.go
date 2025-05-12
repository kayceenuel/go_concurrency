package cache

import (
	"sync/atomic"
)

// Statistics tracks cache usage statistics
type Statistics struct {
	Reads              int64   // Total number of read operations
	Writes             int64   // Total number of write operations
	Hits               int64   // Number of cache hits
	Misses             int64   // Number of cache misses
	Evictions          int64   // Number of entries evicted
	NeverReadCount     int64   //Total evicted items that were never read
	CurrentNeverRead   int     //Current items never read (calculated on demand)
	AverageAccessCount float64 // Average access  count (calculated on demand)

}

// newStatistics creates a new statistics tracker
func newStatistics() *Statistics {
	return &Statistics{}
}

// IncrementReads increments the read counter
func (s *Statistics) IncrementReads() {
	atomic.AddInt64(&s.Reads, 1)
}

// IncrementWrites increments the writes counter
func (s *Statistics) IncrementWrites() {
	atomic.AddInt64(&s.Writes, 1)
}

// IncrementHits increments the hits counter
func (s *Statistics) IncrementHits() {
	atomic.AddInt64(&s.Hits, 1)
}

// IncrementMisses increments the misses counter
func (s *Statistics) IncrementMisses() {
	atomic.AddInt64(&s.Misses, 1)
}

// IncrementEvictions increments the evictions counter
func (s *Statistics) IncrementEvictions() {
	atomic.AddInt64(&s.Evictions, 1)
}

// IncrementNeverRead increments the counter for evicted items that were never read
func (s *Statistics) IncrementNeverRead() {
	atomic.AddInt64(&s.NeverReadCount, 1)
}

// GetHitReate calculates the cache hit rate
func (s *Statistics) GetHitRate() float64 {
	reads := atomic.LoadInt64(&s.Reads)
	if reads == 0 {
		return 0
	}
	hits := atomic.LoadInt64(&s.Hits)
	return float64(hits) / float64(reads)
}
