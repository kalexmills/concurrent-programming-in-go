package main

import (
	"fmt"
	"sync"
)

// Based on: https://www.youtube.com/watch?v=KBZlN0izeiY

func main() {

	var wg sync.WaitGroup
	bc := NewBufferedChannel(4)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			bc.Send(i)
			fmt.Printf("sent %d\n", i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			x := bc.Receive()
			fmt.Printf("received %d\n", x)
		}
	}()

	wg.Wait()
}

type BufferedChannel struct {
	buffer       []int
	head         int // head is the index of the next empty slot (if there is one)
	tail         int // tail is the index of the first unread item
	isEmpty      bool
	mut          *sync.Mutex
	emptyWaiting *sync.Cond
	fullWaiting  *sync.Cond
}

func NewBufferedChannel(size int) *BufferedChannel {
	mut := &sync.Mutex{}
	return &BufferedChannel{
		buffer:       make([]int, size),
		isEmpty:      true,
		mut:          mut,
		emptyWaiting: sync.NewCond(mut),
		fullWaiting:  sync.NewCond(mut),
	}
}

func (bc *BufferedChannel) Send(item int) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && !bc.isEmpty {
		bc.fullWaiting.Wait()
	}

	bc.buffer[bc.head] = item
	bc.head = (bc.head + 1) % len(bc.buffer)
	bc.isEmpty = false
	bc.emptyWaiting.Signal()
}

func (bc *BufferedChannel) Receive() int {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && bc.isEmpty {
		bc.emptyWaiting.Wait()
	}

	item := bc.buffer[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.buffer)

	if bc.head == bc.tail {
		bc.isEmpty = true
	}
	bc.fullWaiting.Signal()

	return item
}
