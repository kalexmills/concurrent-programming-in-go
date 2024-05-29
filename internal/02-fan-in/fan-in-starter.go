package main

import (
	"fmt"
	"sync"
)

// https://go.dev/play/p/WBkddqGE9T_D

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	n := 10
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				if i == 0 { // first producer is responsible for closing
					wg.Wait()
					close(ch)
				}
			}()
			for j := 0; j < 10; j++ {
				ch <- j
			}
			fmt.Printf("producer %d done\n", i)
		}()
	}

	go func() {
		for data := range ch {
			fmt.Println("data received:", data)
		}
		fmt.Println("consumer done")
	}()

	wg.Wait()
	fmt.Println("main func done")
}
