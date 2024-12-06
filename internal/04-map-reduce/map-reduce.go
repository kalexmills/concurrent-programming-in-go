package main

import (
	"fmt"
	"strings"
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
	}()

	// start mappers
	for i := 0; i < numMappers; i++ {
		go func() {
			for line := range linesChan {
				words := strings.Split(line, " ")
				for _, word := range words {
					lowered := strings.ToLower(words[i])
					reducerID := int(lowered[0]-'a') % numReducers
					wordsChan[reducerID] <- word
				}
			}
		}()
	}

	// start reducers
	for i := 0; i < numReducers; i++ {
		go func() {
			localCount := make(map[string]int)
			for word := range wordsChan[i] {
				localCount[word]++
			}
			countChan <- localCount
		}()
	}

	// start consumer
	go func() {
		for count := range countChan {
			fmt.Printf("received count: %v\n", count)
		}
	}()
}
