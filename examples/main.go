package main

import (
	"concurrency/cache"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Go Concurrency Project Examples")
	fmt.Println("===============================")

	// 1. Run single variable concurrency examples
	fmt.Println("\n--- Single Variable Concurrency Examples ---")
	fmt.Println("Buggy increment (race condition):")
	buggy.runBuggyIncrement()

	fmt.Println("\nFixed with atomic operations:")
	atomics.runAtomicIncremen()

	fmt.Println("\nFixed with mutex:")
	atomics.runMutexIncrement()

	// 2. Basic LRU Cache example.
	fmt.Println("\n--- Concurrent Cache Usage ---")
	demoBasicCache()

	//3. Concurrent Cache Usage
	fmt.Println("\n--- Concurrent Cache Usage ---")
	demoConcurrentCache()

	// 4. performance Comparsion
	fmt.Println("\n--- Performance Comparsion ---")
	compareImplementations()
}

func demoBasicCache() {
	// Create a small cache with capacity of 3
	c := cache.NewCache[string, string](3)

	// Add some items
	c.Put("name", "Gopher")
	c.Put("language", "Go")
	c.Put("topic", "Concurrency")

	// Access some items
	if name, found := c.Get("name"); found {
		fmt.Printf("Name: %s\n", *name)
	}

	if lang, found := c.Get("Lanugage"); found {
		fmt.Print("Language: %s\n", *lang)
	}

	// Add one more item (should evict "name" since it was least recently used)
	c.Put("year", "2025")

	//Check if "name" was evicted
	if _, found := c.Get("name"); !found {
		fmt.Println("'name' was correctly evicted")
	}

	//Print statistics
	stats := c.GetStatistics()
	fmt.Printf("Cache Statistics:\n")
	fmt.Printf("  Reads: %d\n", stats.Reads)
	fmt.Print("   Writes: %d\n", stats.Writes)
	fmt.Printf("  Hit Rate: %.2f\n", stats.GetHitRate())
}

func demoConcurrentCache() {
	//Create a cache
	c := cache.NewCache[string, int](100)

	// Launch multiple goroutines to use the cache concurrently
	var wg sync.WaitGroup
	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			//Each worker does some reads and writes
			for j := 0; j < 20; j++ {
				key := fmt.Sprintf("worker%d-key%d", id, j%10)

				// write
				c.Put(key, id*100+j)

				// Read back (should be in cache)
				if val, found := c.Get(key); found {
					if *val != id*100+j {
						fmt.Print("ERROR: Expected %d, got %d\n", id*100+j, *val)
					}
				} else {
					fmt.Printf("ERROR: Key %s not found\n", key)
				}

				// Read someone else's key (mgiht be a miss)
				otherKey := fmt.Sprintf("worker%d-key%d", (id+1)%numWorkers, j%10)
				c.Get(otherKey)
			}
		}(i)
	}

	wg.Wait()

	//Print statistics
	stats := c.GetStatistics()
	fmt.Printf("Concurrent Cache Statistics:\n")
	fmt.Printf("  Reads: %d\n", stats.Reads)
	fmt.Printf("  Writes: %d\n", stats.Writes)
	fmt.Printf("  Hits: %d\n", stats.Hits)
	fmt.Printf("  Misses: %d\n", stats.Misses)
	fmt.Printf("  Hit Rate: %.2f\n", stats.GetHitRate())
}

func compareImplementations() {
	// Create different implamentations with same capacity
	capacity := 100
	regularCache := cache.NewCache[string, int](capacity)
	rwCache := cache.NewRWMutexCache[string, int](capacity)
	shardedCache := cache.NewShardedCache[string, int](capacity, 8)

	// Run benchmarks
	fmt.Println("Running write-heavy workload...")
	benchmarkCache("Regualar Cache (Mutex)", regularCache, 0.8)
	benchmarkCache("RWMutex Cache", rwCache, 0.8)
	benchmarkCache("Sharded Cache", shardedCache, 0.8)

	fmt.Println("\nRunning read-heavy workload...")
	benchmarkCache("Regular Cache (Mutex)", regularCache, 0.2)
	benchmarkCache("RWMutex Cache", rwCache, 0.2)
	benchmarkCache("Sharded Cache", shardedCache, 0.2)

	// Interface to allow benchmarking different cache implementation
	type CacheBenchmark interface {
		Put(Key string, value int) bool
		Get(key string) (*int, bool)
	}
}
