package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// contexts carry timeouts and deadlines
	// have a Done() channel which is closed when timeout or deadline has elapsed.

	ch := make(chan int) // work queue
	n := 100

	var wg sync.WaitGroup
	for id := 0; id < 10; id++ {
		wg.Add(1)
		go func() {
			defer func() {
				fmt.Printf("worker %d done\n", id)
				wg.Done()
			}()
			// blocks until next value in channel is ready
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						fmt.Printf("worker %d done; channel closed\n", id)
						return
					}
					DoRPC(ctx, id, msg)
				case <-ctx.Done():
					fmt.Printf("worker %d done; context cancelled\n", id)
					return
				}
			}
		}()
	}

loop:
	for i := 0; i < n; i++ {
		// context-aware sending to a channel.
		select {
		case ch <- i:
			// do nothing in response
		case <-ctx.Done():
			fmt.Println("sender context cancelled")
			break loop
		}
	}
	close(ch)
	wg.Wait()
	fmt.Println("end of main")
}

func DoRPC(ctx context.Context, workerID int, val int) {
	fmt.Printf("worker: %d, executing RPC with data: %v\n", workerID, val)
	time.Sleep(150 * time.Millisecond)
	// pass context down into the HTTP library.
}
