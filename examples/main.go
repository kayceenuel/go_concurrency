package main

import (
	"concurrency/cache"
	"fmt"
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
}
