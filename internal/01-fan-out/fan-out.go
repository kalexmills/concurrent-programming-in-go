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
			defer wg.Done()
			DoRPC(i)
		}()
	}
	// at some point in future: go runtime will start 10 goroutines
	//    running DoRPC
	wg.Wait()
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("doing an RPC w/ data: %d\n", data)
}
