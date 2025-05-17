package cache

import (
	"fmt"
	"sync"
	"testing"
)

// TestConcurrentAccess verifies the cache is safe for concurrent use
func TestConcurrentAccess(t *testing.T) {
	// Create a new cache instance
	cache := NewCache[string, int](100)

	var wg sync.WaitGroup
	numWorkers := 10
	opsPerWorker := 100

	// Launch multiple goroutines to perfore operations concurrently
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < opsPerWorker; j++ {
				// Alternate between reads and writes
				if j%2 == 0 {
					key := fmt.Sprintf("key%d", j%100)
					cache.Put(key, workerID*j)
				} else {
					key := fmt.Sprintf("key%d", (j-1)%100)
					cache.Get(key)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify we have sensible the statistics
	stats := cache.GetStatistics()
	expectedOps := numWorkers * opsPerWorker
	if int(stats.Reads)+int(stats.Writes) != expectedOps {
		t.Errorf("Expected %d total operations, got %d reads and %d writes",
			expectedOps, stats.Reads, stats.Writes)
	}
}

// TestConcurrentRWMUtexCache tests the RWMutex implementation
func TestConcurrentRWMutexCache(t *testing.T) {
	cache := NewRWMutexCache[string, int](100)

	var wg sync.WaitGroup
	numReaders := 8
	numWriters := 2
	opsPerWorker := 1000

	// Launch reader  goroutines (mostly reads)
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("key%d", j%100)

				//Mostly reads, occasional writes
				if j%10 == 0 {
					cache.Put(key, workerID*j)
				} else {
					cache.Get(key)
				}
			}
		}(i)
	}

	// launch writer goroutines (mostly writes)
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("key%d", j%100)

				// Mostly writes, occastinal reads
				if j%10 == 0 {
					cache.Get(key)
				} else {
					cache.Put(key, workerID*j)
				}
			}
		}(i + numReaders)
	}

	wg.Wait()

	// Verify we have sensible statistics
	stats := cache.GetStatistics()
	expectedOps := (numReaders + numWriters) * opsPerWorker
	if int(stats.Reads)+int(stats.Writes) != expectedOps {
		t.Errorf("Expected %d total operations, got %d reads and %d writes",
			expectedOps, stats.Reads, stats.Writes)
	}
}
