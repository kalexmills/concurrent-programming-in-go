package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	ctx.Done() // a special channel that's closed when the context is cancelled.

	numMsgs := 100
	numWorkers := 10
	var wg sync.WaitGroup

	workQueue := make(chan int)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case msgID, ok := <-workQueue:
					if !ok {
						fmt.Printf("worker %d found channel closed\n", workerID)
						return
					}
					SlowRPCs(ctx, workerID, msgID)
				case <-ctx.Done():
					fmt.Printf("worker %d cancelled\n", workerID)
					return
				}
			}
		}(i)
	}

sendLoop:
	for msgID := 0; msgID < numMsgs; msgID++ {
		select {
		case workQueue <- msgID: // blocks
			// run some code if this channel op happens first
		case <-ctx.Done():
			break sendLoop
		}
	}
	close(workQueue) // it is always the senders responsibility to close the channel

	wg.Wait()
}

func SlowRPCs(ctx context.Context, workerID, msgID int) {
	time.Sleep(500 * time.Millisecond) // fake latency.
	fmt.Printf("worker #%d sending RPC with value %d\n", workerID, msgID)
	// TODO: keep passing ctx down into the network library
}
