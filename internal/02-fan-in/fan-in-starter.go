package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)

	var producerWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		producerWg.Add(1)
		go func(id int) {
			defer func() {
				producerWg.Done()
				if id == 0 {
					producerWg.Wait()
					close(ch)
				}
			}()
			for i := 0; i < 10; i++ {
				ch <- i
			}
			fmt.Printf("producer %d stopped\n", id)
		}(i)
	}

	// Each producer should generate and send 10 integers to the consumer using ch.
	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for i := range ch {
			fmt.Printf("received %d\n", i)
		}
		fmt.Println("consumer stopped")
	}()
	// TODO: start 1 consumer goroutine
	// The consumer should receive all of the integers from ch and print them out.

	// Challenge: Use wait-groups to ensure that every goroutine returns before the main() func stops.

	consumerWg.Wait()
	fmt.Println("main goroutine complete")
}
