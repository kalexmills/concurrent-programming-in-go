package main

import (
	"fmt"
	"sync"
)

// https://go.dev/play/p/0YBVXu1N2CR

func main() {
	fanInChan := make(chan int)
	var producerWg sync.WaitGroup

	for i := 0; i < 10; i++ {
		producerWg.Add(1)
		go func(workerID int) {
			defer func() {
				producerWg.Done()
				if workerID == 0 {
					producerWg.Wait()
					close(fanInChan)
				}
			}()
			for j := 0; j < 10; j++ {
				fanInChan <- j
			}
			fmt.Printf("producer %d finished\n", workerID)
		}(i)
	}

	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for msg := range fanInChan {
			fmt.Printf("received %d\n", msg)
		}
		fmt.Println("consumer finished")
	}()

	producerWg.Wait()
	consumerWg.Wait()
}
