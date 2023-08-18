package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	n := 10

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(msgID int) {
			defer wg.Done()
			DoRPC(msgID)
		}(i)
		// guarantee: goroutine will start
	}
	wg.Wait() // blocks the main goroutine until the counter is zero
	fmt.Println("program ended")
}

// DoRPC fakes a remote procedure call.
func DoRPC(msgID int) {
	fmt.Println("sending message id:", msgID)
	time.Sleep(100 * time.Millisecond) // fake RPC.
	// blocking call -- stops the currently running goroutine.
	fmt.Println("message", msgID, "was sent")
}
