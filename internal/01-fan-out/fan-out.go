package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	workerPool := make(chan int)

	var wg sync.WaitGroup
	for id := 0; id < 10; id++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range workerPool { // loop until channel is closed
				DoRPC(data)
			}
			fmt.Printf("receiver %d is done!\n", id)
		}()
	}
	// at some point in future: go runtime will start 10 goroutines
	//    running DoRPC

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			workerPool <- i // block until a receiver is available (when channel is full)
		}
		close(workerPool)
		fmt.Println("sender is done!")
	}()

	wg.Wait()
	fmt.Println("main is done!")
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("doing an RPC w/ data: %d\n", data)
}
