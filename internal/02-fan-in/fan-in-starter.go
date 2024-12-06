package main

import (
	"fmt"
	"sync"
)

// https://goplay.tools/snippet/IVdAC39Drkx

func main() {
	ch := make(chan int)
	var cwg sync.WaitGroup
	var pwg sync.WaitGroup

	for i := 0; i < 10; i++ {
		pwg.Add(1)
		go func() {
			defer func() {
				pwg.Done()
				if i == 0 {
					pwg.Wait()
					close(ch)
				}
			}()
			for j := 0; j < 10; j++ {
				ch <- i * j
			}
			fmt.Printf("producer %d shut down\n", i)
		}()
	}

	cwg.Add(1)
	go func() {
		defer cwg.Done()
		for data := range ch {
			fmt.Printf("received %d\n", data)
		}
		fmt.Println("consumer shut down")
	}()

	pwg.Wait()
	cwg.Wait()
}
