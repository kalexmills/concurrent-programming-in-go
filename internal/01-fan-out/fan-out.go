package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(msgID int) {
			DoRPC(msgID)
			wg.Done()
		}(i)
		// ask the runtime to start a new goroutine
		// at some point in the future a new goroutine will start
		// and run DoRPC.
	}
	// i = 10
	// guaranteed that 10 goroutines will start sometime...
	// wg count cannot become 0 until every goroutine has stopped.
	wg.Wait()
}

func DoRPC(msgID int) {
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("sending message %d\n", msgID)
}
