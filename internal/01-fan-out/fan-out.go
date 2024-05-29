package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			LongRunningRPC(i)
			wg.Done()
		}()
		// guarantee: at some point in future, LongRunningRPC will run in parallel.
	}
	wg.Wait()
}

func LongRunningRPC(data int) {
	fmt.Printf("sending %d to server\n", data)
	time.Sleep(150 * time.Millisecond)
}
