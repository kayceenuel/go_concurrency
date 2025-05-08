package exercises

import (
	"fmt"
	"sync"
)

// RunBuggyIncrement demonstrates the problem with unsynchronized access
func RunBuggyIncrement() {
	for i := 0; i < 5; i++ {
		buggyIncrement()
	}
}

func buggyIncrement() {
	var x = 0
	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			x = x + 1 // This operation is not atomic!
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println("Buggy increment result:", x)
	// Almost certainly won't be 1000 due to race conditions
}
