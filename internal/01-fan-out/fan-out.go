package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel() // for resource cleanup

	var wg sync.WaitGroup
	numWorkers := 10

	workQueue := make(chan int)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) { // worker goroutines
			defer wg.Done() // -- 'happens before' line 33
			for {           // blocks until a msg is received
				select {
				case msg, ok := <-workQueue:
					if !ok {
						fmt.Println("channel closed! from worker", workerID)
						return
					}
					DoRPC(ctx, workerID, msg)
				case <-ctx.Done():
					fmt.Println("worker", workerID, "ended") // (happens after close at line 32)
					return
				}
			}
		}(i)
	}

loop:
	for i := 0; i < 100; i++ {
		// racing two channel operations against one another
		select {
		case workQueue <- i: // blocks until a receiver is available
			// run this code
		case <-ctx.Done():
			fmt.Printf("sender was cancelled while sending message %d\n", i)
			break loop
		}
	}
	close(workQueue) // closes the channel (happens before line 24)
	wg.Wait()        // blocks until the counter is zero (i.e. until all goroutines have finished)
	fmt.Println("program ended")
}

// DoRPC fakes a remote procedure call.
func DoRPC(ctx context.Context, workerID int, msg int) {
	fmt.Printf("sending message %d from worker %d\n", msg, workerID)
	// TODO: use the ctx in the real RPC.
	time.Sleep(100 * time.Millisecond) // fake RPC.
	// blocking call -- stops the currently running goroutine.
	fmt.Println("worker", workerID, ": message", msg, "was sent")
}
