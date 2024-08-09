package main

import (
	"fmt"
	"strings"
	"sync"
)

// https://go.dev/play/p/CU8lt4mIfl

var lines = []string{
	"Lorem Ipsum is simply dummy text of the printing and typesetting industry.",
	"Lorem Ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book.",
	"It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
}

func main() {
	linesCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, line := range lines {
			linesCh <- line
		}
		close(linesCh) // happens-before end of mapper range-over-channel
		fmt.Println("producer done")
	}()

	numMappers := 3
	numReducers := 3
	wordsChs := make([]chan string, numReducers)
	for i := 0; i < numReducers; i++ {
		wordsChs[i] = make(chan string)
	}

	for i := 0; i < numMappers; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				if i == 0 {
					wg.Wait()
					for _, ch := range wordsChs {
						close(ch) // happens before end of reducers range-over-channel
					}
				}
			}()

			// mapper's range-over-channel
			for line := range linesCh { // end happens before close of wordsChs
				words := strings.Split(strings.ToLower(line), " ")
				for _, word := range words {
					key := int(word[0]-'a') % numReducers
					wordsChs[key] <- word
				}
			}
			fmt.Println("mapper done:", i)
		}()
	}

	var reducerWg sync.WaitGroup
	countCh := make(chan map[string]int)
	for i := 0; i < numReducers; i++ {
		localCount := make(map[string]int)
		reducerWg.Add(1)
		go func() {
			defer func() {
				reducerWg.Done()
				if i == 0 {
					reducerWg.Wait()
					close(countCh) // happens-before end of consumer's range-over-channel
				}
			}()
			// reducers range-over-channel
			for word := range wordsChs[i] { // end happens before close of countCh
				localCount[word]++
			}
			countCh <- localCount
			fmt.Println("reducer done:", i)
		}()
	}

	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()      // happens-before end of consumerWg.Wait() at line 96
		for count := range countCh { // consumer's range-over-channel
			fmt.Println("got count:", count)
		}
		fmt.Println("consumer done")
	}()

	consumerWg.Wait()
	fmt.Print("main done")
}
