package buggy

import (
	"fmt"
	"sync"
)

func main() {
	// Run the buggy increment multiple times to show inconsistency
	fmt.Println("Running buggy increment example 5 times:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Run %d: ", i+1)
		RunBuggyIncrement()
	}
}

func RunBuggyIncrement() {
	var x = 0
	var wg sync.WaitGroup

	// Start 1000 goroutines that will each try to increment x
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x = x + 1 // This operation is not atomic!
			wg.Done()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print the final value - it should be 1000 if all increments were performed correctly
	fmt.Println("final value of x:", x)
	// Almost certainly won't be 1000 due to race conditions
}
