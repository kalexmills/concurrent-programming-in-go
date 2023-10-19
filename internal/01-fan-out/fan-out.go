package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	numMsgs := 100
	numWorkers := 10
	var wg sync.WaitGroup

	workQueue := make(chan int)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for msgID := range workQueue {
				SlowRPCs(workerID, msgID)
			} // only way to break from a channel range loop is when channel is closed
		}(i)
	}

	for msgID := 0; msgID < numMsgs; msgID++ {
		workQueue <- msgID // send the ith message (could be a blocking operation)
	}
	close(workQueue) // it is always the senders responsibility to close the channel

	wg.Wait()
}

func SlowRPCs(workerID, msgID int) {
	time.Sleep(500 * time.Millisecond) // fake latency.
	fmt.Printf("worker #%d sending RPC with value %d\n", workerID, msgID)
}
