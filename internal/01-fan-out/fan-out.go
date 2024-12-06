package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			DoRPC(i)
		}()
		// guarantee: at some point in future DoRPC will start (on another goroutine).
	}
	wg.Wait() // guaranteed: 10 goroutines will start later
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("sent data %d\n", data)
}
