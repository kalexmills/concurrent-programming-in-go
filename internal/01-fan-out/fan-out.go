package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
	defer cancel()

	var wg sync.WaitGroup
	ch := make(chan int)

	// block until timeout has elapsed
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-ch:
					if !ok { // channel was closed
						fmt.Printf("goroutine %d complete; channel closed\n", i)
						return
					}
					DoRPC(data)
				case <-ctx.Done():
					fmt.Printf("goroutine %d complete; timeout\n", i)
					return
				}
			}
		}()
	}

loop:
	for i := 0; i < 100; i++ {
		select {
		case ch <- i:
			// do nothing
		case <-ctx.Done():
			fmt.Println("timeout elapsed; sender")
			break loop
		}
	}
	close(ch) // signal that no further data is coming

	wg.Wait()
	fmt.Println("main thread complete")
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("sent data %d\n", data)
}
