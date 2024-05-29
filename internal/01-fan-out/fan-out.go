package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer done()
	ch := make(chan int)
	n := 10

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-ch:
					if !ok {
						fmt.Printf("worker #%d stopped; channel closed\n", i)
						return
					}
					LongRunningRPC(ctx, data)
				case <-ctx.Done():
					fmt.Printf("worker #%d stopped; deadline expired\n", i)
					return
				}
			}
		}()
	}

loop:
	for msg := 0; msg < 100; msg++ {
		select {
		case ch <- msg:
		case <-ctx.Done():
			break loop
		}
	}

	close(ch)
	wg.Wait()
	fmt.Println("ending main func")
}

func LongRunningRPC(ctx context.Context, data int) {
	fmt.Printf("sending %d to server\n", data)
	time.Sleep(150 * time.Millisecond) // TODO: pass context to networking library
}
