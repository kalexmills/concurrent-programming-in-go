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
	defer cancel()

	var wg sync.WaitGroup
	n := 100

	ch := make(chan int)

	wg.Add(10)
	for i := 0; i < 10; i++ { // start 10 consumers
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case msgID, ok := <-ch:
					if !ok {
						fmt.Printf("channel closed, worker %d stopped\n", workerID)
						return
					}
					DoRPC(ctx, workerID, msgID) // blocking call (150ms)
				case <-ctx.Done():
					fmt.Printf("worker %d stopped\n", workerID)
					return
				}
			}
		}(i)
	}
	// guarantee: at some point in future 10 workers will start on separate
	//            goroutines

loop:
	for i := 0; i < n; i++ {
		select {
		case ch <- i: // sending blocks until a receiver is available
			fmt.Printf("sent message %d\n", i)
		case <-ctx.Done(): // wait until 250ms have elapsed
			break loop
		}
	}
	close(ch) // closes the channel -- no more data to send!
	wg.Wait() // wait for all the goroutines to halt
	fmt.Println("main goroutine done!")
}

func DoRPC(ctx context.Context, workerID, msgID int) {
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("worker %d: sending message %d\n", workerID, msgID)
	// TODO: actually use the context in real networking code
}
