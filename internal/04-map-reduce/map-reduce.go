package main

import (
	"fmt"
	"strings"
	"sync"
)

// https://go.dev/play/p/CU8lt4mIflo

var lines = []string{
	"Lorem Ipsum is simply dummy text of the printing and typesetting industry",
	"Lorem Ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book",
	"It has survived not only five centuries but also the leap into electronic typesetting remaining essentially unchanged It was popularised in the with the release of Letraset sheets containing Lorem Ipsum passages and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum",
}

func main() {
	linesChan := make(chan string)
	wordsChan := make([]chan string, 26)
	for i := 0; i < 26; i++ {
		wordsChan[i] = make(chan string)
	}
	wordCountChan := make(chan map[string]int)
	numMappers := 3
	numReducers := 26

	// start the producer to feed lines into the system
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, line := range lines {
			linesChan <- strings.ToLower(line)
		}
		close(linesChan)
		fmt.Println("producer complete")
	}()

	// start the mappers to split lines into words and send to the correct reducers
	var wgMapper sync.WaitGroup
	for i := 0; i < numMappers; i++ {
		wgMapper.Add(1)
		go func(id int) {
			defer func() {
				wgMapper.Done()
				if id == 0 {
					wgMapper.Wait()
					for i := 0; i < len(wordsChan); i++ {
						close(wordsChan[i])
					}
				}
				fmt.Printf("mapper %d complete\n", id)
			}()
			for line := range linesChan {
				tokens := strings.Split(line, " ")
				for _, token := range tokens {
					token = strings.TrimSpace(token)
					idx := int(token[0] - 'a')
					wordsChan[idx] <- token
				}
			}
		}(i)
	}

	// start the reducers to obtain accurate wordcounts
	var wgReducer sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		wgReducer.Add(1)
		go func(workerID int) {
			defer func() {
				wgReducer.Done()
				if workerID == 0 {
					wgReducer.Wait()
					close(wordCountChan)
				}
				fmt.Printf("reducer %d complete\n", workerID)
			}()

			count := make(map[string]int)
			for word := range wordsChan[workerID] {
				count[word]++
			}
			wordCountChan <- count
		}(i)
	}

	// start the consumer to receive and print out counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		for counts := range wordCountChan {
			fmt.Printf("counts received: %v\n", counts)
		}
		fmt.Println("consumer complete!")
	}()

	wg.Wait()
	fmt.Println("main complete")
}
