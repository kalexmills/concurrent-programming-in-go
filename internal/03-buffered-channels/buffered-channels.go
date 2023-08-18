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
			err := ch.Send(i)
			if err != nil {
				fmt.Println("sender", id, ": error received:", err)
			} else {
				fmt.Println("sender", id, ": value sent:", i)
			}
		}
		wg.Done()
	}(0)

	go func(id int) {
		for i := 0; i < 10; i++ {
			val, err := ch.Receive()
			if err != nil {
				fmt.Println("receiver", id, "error received:", err)
			} else {
				fmt.Println("receiver", id, "value received:", val)
			}
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

	mut *sync.Mutex
}

func NewBufferedChannel(size int) *BufferedChannel {
	return &BufferedChannel{
		buffer: make([]int, size),
		mut:    &sync.Mutex{},
	}
}

var ErrEmpty = errors.New("channel was empty")
var ErrFull = errors.New("channel was full")

func (bc *BufferedChannel) Send(val int) error {
	bc.mut.Lock()
	defer bc.mut.Unlock()
	// "critical section"

	if bc.head == bc.tail && bc.isFull {
		return ErrFull
	}

	bc.buffer[bc.head] = val
	bc.head = (bc.head + 1) % len(bc.buffer)

	if bc.head == bc.tail {
		bc.isFull = true
	}

	return nil
}

func (bc *BufferedChannel) Receive() (int, error) {
	bc.mut.Lock()
	defer bc.mut.Unlock()
	// begin "critical section"

	if bc.head == bc.tail && !bc.isFull { // TODO: fill in this conditional
		return 0, ErrEmpty
	}

	val := bc.buffer[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.buffer)

	bc.isFull = false

	return val, nil
}
