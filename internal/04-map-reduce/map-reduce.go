package main

import (
	"fmt"
	"strings"
	"sync"
)

var lines = []string{
	"Lorem Ipsum is simply dummy text of the printing and typesetting industry.",
	"Lorem Ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book.",
	"It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
}

func main() {
	numMappers := 3
	numReducers := 5
	linesChan := make(chan string)
	wordsChan := make([]chan string, numReducers)
	countChan := make(chan map[string]int)
	for i := 0; i < numReducers; i++ {
		wordsChan[i] = make(chan string)
	}

	// start producer
	go func() {
		for _, line := range lines {
			linesChan <- line
		}
		close(linesChan)
		fmt.Println("producer finished")
	}()

	// start mappers
	var mwg sync.WaitGroup
	for i := 0; i < numMappers; i++ {
		mwg.Add(1)
		go func() {
			defer func() {
				mwg.Done()
				if i == 0 {
					mwg.Wait()
					for _, ch := range wordsChan {
						close(ch)
					}
				}
			}()
			for line := range linesChan {
				words := strings.Split(line, " ")
				for _, word := range words {
					lowered := strings.ToLower(word)
					reducerID := int(lowered[0]-'a') % numReducers
					wordsChan[reducerID] <- lowered
				}
			}
			fmt.Printf("mapper %d finished\n", i)
		}()
	}

	// start reducers
	var rwg sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		rwg.Add(1)
		go func() {
			defer func() {
				rwg.Done()
				if i == 0 {
					rwg.Wait()
					close(countChan)
				}
			}()
			localCount := make(map[string]int)
			for word := range wordsChan[i] {
				localCount[word]++
			}
			countChan <- localCount
			fmt.Printf("reducer %d finished\n", i)
		}()
	}

	// start consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for count := range countChan {
			fmt.Printf("received count: %v\n", count)
		}
		fmt.Println("consumer finished")
	}()

	wg.Wait()
	fmt.Println("main finished")
}
