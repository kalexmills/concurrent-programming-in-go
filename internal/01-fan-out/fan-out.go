package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	numWorkers := 10

	workQueue := make(chan int) // un

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) { // worker goroutines
			defer wg.Done() // -- 'happens before' line 33
			// keep receiving messages as long as there are messages
			// to send... until the channel is closed
			for msg := range workQueue {
				DoRPC(workerID, msg)
			}
			fmt.Println("worker", workerID, "ended") // (happens after close at line 32)
		}(i)
	}

	for i := 0; i < 100; i++ {
		// blocks until a receive operation happens on another goroutine
		workQueue <- i // sending 100 messages through the channel (happens before line 22)
	}
	close(workQueue) // closes the channel (happens before line 24)
	//wg.Wait()        // blocks the main goroutine until the counter is zero
	fmt.Println("program ended")
}

// DoRPC fakes a remote procedure call.
func DoRPC(workerID int, msg int) {
	fmt.Printf("sending message %d from worker %d\n", msg, workerID)
	time.Sleep(100 * time.Millisecond) // fake RPC.
	// blocking call -- stops the currently running goroutine.
	fmt.Println("worker", workerID, ": message", msg, "was sent")
}
