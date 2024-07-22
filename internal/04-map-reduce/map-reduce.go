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
	numReducers := 5

	wordsChan := make([]chan string, numReducers)
	linesChan := make(chan string)
	countChan := make(chan map[string]int)

	for i := 0; i < numReducers; i++ {
		wordsChan[i] = make(chan string)
	}

	go func() {
		for _, line := range lines {
			linesChan <- line
		}
		close(linesChan)
		fmt.Printf("lines producer done\n")
	}()

	var wgMapper sync.WaitGroup
	numMappers := 3
	for i := 0; i < numMappers; i++ {
		wgMapper.Add(1)
		go func() {
			defer func() {
				wgMapper.Done()
				if i == 0 {
					wgMapper.Wait()
					for _, ch := range wordsChan {
						close(ch)
					}
				}
			}()

			for line := range linesChan {
				words := strings.Split(strings.ToLower(line), " ")
				for _, word := range words {
					key := int(word[0] - 'a')
					wordsChan[key%numReducers] <- word
				}
			}
			fmt.Printf("mapper %d done\n", i)
		}()
	}

	var wgReducer sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		wgReducer.Add(1)
		go func() {
			defer func() {
				wgReducer.Done()
				if i == 0 {
					wgReducer.Wait()
					close(countChan)
				}
			}()

			countMap := make(map[string]int)
			for word := range wordsChan[i] {
				countMap[word] += 1
			}
			countChan <- countMap
			fmt.Printf("reducer %d done\n", i)
		}()
	}

	for countMap := range countChan {
		fmt.Println("received map:", countMap)
	}

	fmt.Println("main func completed")
}
