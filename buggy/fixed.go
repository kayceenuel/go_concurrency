package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x.Add(1)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("final value of x:", x.Load())
}
