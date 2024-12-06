package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan int) // TODO: unbuffered channels sync properties

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// continue receiving data from channel until.... no more data.
			for data := range ch {
				DoRPC(data)
			}
			fmt.Printf("go routine %d complete\n", i)
		}()
	}

	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch) // signal that no further data is coming

	wg.Wait()
	fmt.Println("main thread complete")
}

func DoRPC(data int) {
	time.Sleep(time.Millisecond * 150)
	fmt.Printf("sent data %d\n", data)
}
