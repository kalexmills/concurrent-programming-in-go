package main

import "sync"

// https://go.dev/play/p/0YBVXu1N2CR

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// TODO: start 10 'producer' goroutines
	// Each producer should generate and send 10 integers to the consumer using fanInChan.

	// TODO: start 1 consumer goroutine
	// The consumer should receive all of the integers from fanInChan and print them out.

	// Challenge: Use wait-groups to ensure that every goroutine returns before the main() func stops.

	wg.Wait()
}
