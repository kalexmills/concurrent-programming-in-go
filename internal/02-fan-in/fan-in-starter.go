package main

import (
	"fmt"
	"sync"
)

// https://go.dev/play/p/0YBVXu1N2CR

func main() {
	ch := make(chan int)
	var wg1 sync.WaitGroup

	// producers
	for i := 0; i < 10; i++ {
		wg1.Add(1) // counts # sending goroutines
		go func() {
			for j := 0; j < 10; j++ {
				ch <- j
			}
			wg1.Done()
			if i == 0 {
				wg1.Wait()
				close(ch)
			}
		}()
	}

	// consumer
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for x := range ch {
			fmt.Printf("received: %d\n", x)
		}
	}()

	wg2.Wait() // waiting!
}
