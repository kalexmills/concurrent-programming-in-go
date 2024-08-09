package main

import (
	"fmt"
	"sync"
)

// https://go.dev/play/p/0YBVXu1N2CR

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// Each producer should generate and send 10 integers to the consumer using ch.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				ch <- i
			}
			fmt.Printf("producer %d done!\n", i)
		}()
	}

	// TODO: start 1 consumer goroutine
	// The consumer should receive all of the integers from fanInChan and print them out.
	var wgConsumer sync.WaitGroup
	wgConsumer.Add(1)
	go func() {
		defer wgConsumer.Done()
		for data := range ch {
			fmt.Printf("received %d\n", data)
		}
		fmt.Println("consumer done!")
	}()
	// Challenge: Use wait-groups to ensure that every goroutine returns before the main() func stops.

	wgConsumer.Add(1)
	go func() {
		defer wgConsumer.Done()
		wg.Wait()
		close(ch)
		fmt.Println("closer done!")
	}()

	wgConsumer.Wait()
	fmt.Println("main done!")
}
