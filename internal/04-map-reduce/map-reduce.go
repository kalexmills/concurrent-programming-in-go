package main

import (
	"fmt"
	"strings"
	"sync"
)

// play https://go.dev/play/p/0HQJ81CDwwB

var lines = []string{
	"lorem ipsum is simply dummy text of the printing and typesetting industry",
	"lorem ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book",
	"it has survived not only five centuries but also the leap into electronic typesetting remaining essentially unchanged",
	"it was popularised in the with the release of letraset sheets containing lorem ipsum passages and more recently with desktop publishing software like aldus pagemaker including versions of lorem ipsum",
}

func main() {
	numMappers := 3  // any value
	numReducers := 5 // anything from [1, 26]

	lineChan := make(chan string)
	// 1. pass each line into lineChan (fan out)
	go func() {
		defer close(lineChan)
		for _, line := range lines {
			lineChan <- line
		}
		fmt.Println("line feeder is done")
	}()

	// 2. start mappers which read from lineChan, split lines into words, and send via wordChannels to the appropriate
	//    reducers.
	wordChans := make([]chan string, numReducers)
	for i := 0; i < len(wordChans); i++ {
		wordChans[i] = make(chan string)
	}

	var mwg sync.WaitGroup
	for i := 0; i < numMappers; i++ {
		mwg.Add(1)
		go func(id int) {
			defer func() {
				mwg.Done()
				if id == 0 {
					mwg.Wait()
					for _, ch := range wordChans {
						close(ch)
					}
				}
			}()
			for line := range lineChan { // breaks b/c of the close at line 26
				words := strings.Split(line, " ")
				for _, word := range words {
					if len(word) == 0 {
						continue
					}
					letter := word[0]
					idx := int(letter-'a') % numReducers
					wordChans[idx] <- word
				}
			}
			fmt.Printf("mapper %d is done\n", id)
		}(i)
	}

	// 3. start reducers which read from wordChannels, form a local wordcount, and send the results to countChannel.

	countChan := make(chan map[string]int)
	var rwg sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		rwg.Add(1)
		go func(id int) {
			defer func() {
				rwg.Done()
				if id == 0 {
					rwg.Wait()
					close(countChan)
				}
			}()
			wordCount := make(map[string]int)
			for word := range wordChans[id] {
				wordCount[word]++
			}
			countChan <- wordCount
			fmt.Printf("reducer %d is done\n", id)
		}(i)
	}

	// Make sure everything shuts down gracefully before the process ends.
	for count := range countChan {
		fmt.Println("count received:", count)
	}

	fmt.Println("main thread is done")
}
