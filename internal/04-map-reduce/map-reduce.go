package main

import (
	"fmt"
	"strings"
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

	go func() {
		for _, line := range lines {
			linesCh <- line
		}
	}()

	numMappers := 3
	numReducers := 3
	wordsChs := make([]chan string, numReducers)
	for i := 0; i < numReducers; i++ {
		wordsChs[i] = make(chan string)
	}

	for i := 0; i < numMappers; i++ {
		go func() {
			for line := range linesCh {
				words := strings.Split(strings.ToLower(line), " ")
				for _, word := range words {
					key := int(word[0]-'a') % numReducers
					wordsChs[key] <- word
				}
			}
		}()
	}

	countCh := make(chan map[string]int)
	for i := 0; i < numReducers; i++ {
		localCount := make(map[string]int)
		go func() {
			for word := range wordsChs[i] {
				localCount[word]++
			}
			countCh <- localCount
		}()
	}

	go func() {
		for count := range countCh {
			fmt.Println("got count:", count)
		}
	}()
}
