package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	ch := NewBufferedChannel(4)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(id int) { // senders
		for i := 0; i < 10; i++ {
			ch.Send(i)
			fmt.Println("sender", id, ": value sent:", i)
		}
		wg.Done()
	}(0)

	go func(id int) {
		for i := 0; i < 10; i++ {
			val := ch.Receive()
			fmt.Println("receiver", id, "value received:", val)

		}
		wg.Done()
	}(0)
	wg.Wait()
}

type BufferedChannel struct {
	buffer []int

	head int // next empty index in the buffer (if it exists)
	tail int // next unread index in the buffer (if one exists)

	isFull bool // true iff the channel is full

	mut          *sync.Mutex
	fullWaiters  *sync.Cond
	emptyWaiters *sync.Cond
}

func NewBufferedChannel(size int) *BufferedChannel {
	result := &BufferedChannel{
		buffer: make([]int, size),
		mut:    &sync.Mutex{},
	}
	result.fullWaiters = sync.NewCond(result.mut)
	result.emptyWaiters = sync.NewCond(result.mut)
	return result
}

var ErrEmpty = errors.New("channel was empty")
var ErrFull = errors.New("channel was full")

func (bc *BufferedChannel) Send(val int) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && bc.isFull {
		bc.fullWaiters.Wait() // wait for the buffer to be not full
	}

	bc.buffer[bc.head] = val
	bc.head = (bc.head + 1) % len(bc.buffer)

	if bc.head == bc.tail {
		bc.isFull = true
	}
	bc.emptyWaiters.Signal()

	return
}

func (bc *BufferedChannel) Receive() int {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && !bc.isFull {
		bc.emptyWaiters.Wait()
	}

	val := bc.buffer[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.buffer)

	bc.isFull = false
	bc.fullWaiters.Signal()

	return val
}
