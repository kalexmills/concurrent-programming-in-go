package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, done := context.WithTimeout(ctx, time.Millisecond*10000)
	defer done()

	workerPool := make(chan int)

	var wg sync.WaitGroup
	for id := 0; id < 10; id++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		loop:
			for {
				select {
				case data, ok := <-workerPool:
					if !ok {
						break loop
					}
					DoRPC(data)
				case <-ctx.Done():
					fmt.Printf("timeout occurred on receiver %d!\n", id)
					break loop
				}
			}
			fmt.Printf("receiver %d is done!\n", id)
		}()
	}
	// at some point in future: go runtime will start 10 goroutines
	//    running DoRPC

	wg.Add(1)
	go func() {
		defer func() {
			close(workerPool)
			wg.Done()
		}()
		for i := 0; i < 100; i++ {
			select {
			case workerPool <- i: // block until a receiver is available (when channel is full)
				// do nothing
			case <-ctx.Done():
				fmt.Println("sender timed out!")
				return
			}
		}

		fmt.Println("sender is done!")
	}()

	wg.Wait()
	fmt.Println("main is done!")
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("doing an RPC w/ data: %d\n", data)
}
