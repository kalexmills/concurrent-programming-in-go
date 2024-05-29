package main

import (
	"fmt"
	"strings"
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
	countChan := make(chan map[string]int)

	for i := 0; i < 26; i++ {
		wordsChan[i] = make(chan string)
	}

	go func() {
		for _, line := range lines {
			linesChan <- line
		}
	}()

	numMappers := 3
	for i := 0; i < numMappers; i++ {
		go func() {
			for line := range linesChan {
				words := strings.Split(strings.ToLower(line), " ")
				for _, word := range words {
					key := int(word[0] - 'a')
					wordsChan[key] <- word
				}
			}
		}()
	}

	numReducers := 26
	for i := 0; i < numReducers; i++ {
		go func() {
			countMap := make(map[string]int)
			for word := range wordsChan[i] { // TODO: deadlock here
				countMap[word] += 1
			}
			countChan <- countMap
		}()
	}

	for countMap := range countChan {
		fmt.Println("received map:", countMap)
	}
}
