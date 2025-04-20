package main

import (
	"fmt"
	"sync"
)

func main() {
	var x = 0
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x = x + 1 // Buggy: race condition
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("final value of x", x)
}
