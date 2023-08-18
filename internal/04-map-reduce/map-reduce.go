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
	numReducers := 6 // max at 26 -- (because of our hacky choice of hash function)

	lineChan := make(chan string)
	wordChannels := make([]chan string, numReducers)
	for i := 0; i < numReducers; i++ {
		wordChannels[i] = make(chan string)
	}
	countChannel := make(chan map[string]int)

	// mappers
	var mapperWg sync.WaitGroup
	for i := 0; i < numMappers; i++ {
		mapperWg.Add(1)
		go func(id int) {
			defer func() {
				mapperWg.Done()
				if id == 0 { // mapper with id = 0 will close all reducer channels
					mapperWg.Wait() // wait for all mappers to conclude sending
					for i := 0; i < numReducers; i++ {
						close(wordChannels[i]) // close reducer channels.
					}
				}
			}()

			for line := range lineChan {
				// take the first letter in the word and use it to send
				// to the correct reducer
				line = strings.ToLower(line)
				words := strings.Split(line, " ")
				for _, word := range words {
					idx := (int(word[0] - 'a')) % numReducers // dirty trick
					wordChannels[idx] <- word
				}
			}
			fmt.Printf("mapper %d finished\n", id)
		}(i)
	}

	// reducers
	var reducerWg sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		reducerWg.Add(1)
		go func(id int) {
			defer func() {
				reducerWg.Done()
				if id == 0 {
					reducerWg.Wait()
					close(countChannel)
					fmt.Println("count channel closing")
				}
			}()
			// counting all the words seen
			localMap := make(map[string]int)
			for word := range wordChannels[id] {
				localMap[word]++
			}
			countChannel <- localMap
			fmt.Printf("reducer %d finished\n", id)
		}(i)
	}

	// consumer
	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for counts := range countChannel {
			fmt.Println("consumer received: ", counts)
		}
		fmt.Println("consumer done")
	}()

	// feed the mappers each line of the file
	for _, line := range lines {
		lineChan <- line
	}
	close(lineChan)
	fmt.Println("all lines sent!")

	reducerWg.Wait()
	mapperWg.Wait()
	consumerWg.Wait()
}
