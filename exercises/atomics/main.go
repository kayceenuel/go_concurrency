package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("Comparing different increment approaches:")

	// Run the buggy version for comparison
	fmt.Println("\nBuggy increment (with race condition):")
	for i := 0; i < 3; i++ {
		runBuggyIncrement()
	}

	// Run the atomic version
	fmt.Println("\nAtomic increment (thread-safe):")
	for i := 0; i < 3; i++ {
		runAtomicIncrement()
	}

	// Run the mutex version
	fmt.Println("\nMutex increment (thread-safe):")
	for i := 0; i < 3; i++ {
		runMutexIncrement()
	}

	fmt.Println("\nExplanation:")
	fmt.Println("- The buggy version doesn't handle concurrent access properly")
	fmt.Println("- The atomic version uses atomic operations for thread safety")
	fmt.Println("- The mutex version uses a lock to protect the critical section")
}

// Buggy increment - has a race condition
func runBuggyIncrement() {
	var x = 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x = x + 1 // This operation is not atomic!
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("final value of x:", x)
}

// Atomic increment - uses atomic operations
func runAtomicIncrement() {
	var x atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x.Add(1) // This is an atomic operation
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("final value of x:", x.Load())
}

// Mutex increment - uses a mutex
func runMutexIncrement() {
	var x = 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			x = x + 1 // Protected by mutex
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("final value of x:", x)
}
