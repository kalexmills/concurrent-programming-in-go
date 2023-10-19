package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := NewBufferedChannel(4)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch.Send(i)
			fmt.Println("sent value: ", i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			val := ch.Receive()
			fmt.Println("received value: ", val)
		}
	}()
	wg.Wait()
}

//   h
// [ _ _ _ _ _ ]
//   t

type BufferedChannel struct {
	buffer []int
	head   int  // head always points to the next empty slot, if there is one.
	tail   int  // tail always points to the next item to be received (if there is one).
	isFull bool // tells us if the buffer is full (vs. empty)

	mut *sync.Mutex // guards all the variables in a BufferedChannel.

	waitingOnNotFull  *sync.Cond
	waitingOnNotEmpty *sync.Cond
}

func NewBufferedChannel(size int) *BufferedChannel {
	mut := &sync.Mutex{}
	return &BufferedChannel{
		buffer:            make([]int, size),
		mut:               mut,
		waitingOnNotEmpty: sync.NewCond(mut),
		waitingOnNotFull:  sync.NewCond(mut),
	}
}

var ErrBufferFull = errors.New("buffer was full")
var ErrBufferEmpty = errors.New("buffer was empty")

func (ch *BufferedChannel) Send(value int) {
	ch.mut.Lock()
	defer ch.mut.Unlock()

	for ch.isFull {
		ch.waitingOnNotFull.Wait() // releases the lock; lets others do things... and suspends this goroutine
		// we have the lock again!
	}

	ch.buffer[ch.head] = value
	ch.head = (ch.head + 1) % len(ch.buffer)

	if ch.head == ch.tail {
		ch.isFull = true
	}
	ch.waitingOnNotEmpty.Signal()
}

func (ch *BufferedChannel) Receive() int {
	ch.mut.Lock()
	defer ch.mut.Unlock()

	for ch.head == ch.tail && !ch.isFull {
		ch.waitingOnNotEmpty.Wait()
	}

	value := ch.buffer[ch.tail]
	ch.tail = (ch.tail + 1) % len(ch.buffer)

	ch.isFull = false
	ch.waitingOnNotFull.Signal()
	return value
}
