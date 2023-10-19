package main

import (
	"fmt"
	"strings"
)

// Available at: https://go.dev/play/p/aJuViXf6ZYI

var lines = []string{
	"Lorem Ipsum is simply dummy text of the printing and typesetting industry.",
	"Lorem Ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book.",
	"It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
}

func main() {
	numMappers := 3
	numReducers := 26 // max at 26 -- (because of our hacky choice of hash function)

	lineChan := make(chan string)
	wordChannels := make([]chan string, numReducers)
	for i := 0; i < numReducers; i++ {
		wordChannels[i] = make(chan string)
	}
	countChannel := make(chan map[string]int)

	// mappers
	for i := 0; i < numMappers; i++ {
		go func(id int) {
			for line := range lineChan {
				// take the first letter in the word and use it to send
				// to the correct reducer
				line = strings.ToLower(line)
				words := strings.Split(line, " ")
				for _, word := range words {
					idx := int(word[0] - 'a') // dirty trick to find the channel to use
					wordChannels[idx] <- word
				}
			}
			fmt.Printf("mapper %d finished", id)
		}(i)
	}

	// reducers
	for i := 0; i < numReducers; i++ {
		go func(id int) {
			// count all the words seen by this reducer
			localMap := make(map[string]int)
			for word := range wordChannels[id] {
				localMap[word]++
			}
			countChannel <- localMap
			fmt.Printf("reducer %d finished\n", id)
		}(i)
	}

	// consumer
	go func() {
		for counts := range countChannel {
			fmt.Println("consumer received: ", counts)
		}
		fmt.Println("consumer done")
	}()

	// feed the mappers each line of the file
	for _, line := range lines {
		lineChan <- line
	}
	fmt.Println("all lines sent!")

	fmt.Println("main func complete!")
}
