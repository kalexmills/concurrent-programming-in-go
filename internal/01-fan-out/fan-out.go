package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	n := 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(msgID int) {
			SlowRPCs(msgID)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func SlowRPCs(x int) {
	time.Sleep(500 * time.Millisecond) // fake latency.
	fmt.Printf("sending RPC with value %d\n", x)
}
