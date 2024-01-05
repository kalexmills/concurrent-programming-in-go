package main

import (
	"errors"
	"fmt"
	"sync"
)

// Based on: https://www.youtube.com/watch?v=KBZlN0izeiY

func main() {
	ch := NewBufferedChannel(4)

	n := 100

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			ch.Send(i) // blocks until space is available
			fmt.Printf("sent %d\n", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			msg := ch.Receive() // blocks until a msg is available
			fmt.Printf("received %d\n", msg)
		}
	}()

	wg.Wait()
}

type BufferedChannel struct {
	data   []int
	head   int  // next empty space to write to (if one exists)
	tail   int  // next full space to read from
	isFull bool // true when buffer is full

	mut           *sync.Mutex
	untilNotFull  *sync.Cond
	untilNotEmpty *sync.Cond
}

func NewBufferedChannel(size int) *BufferedChannel {
	mut := new(sync.Mutex)
	return &BufferedChannel{
		data:          make([]int, size),
		mut:           mut,
		untilNotEmpty: sync.NewCond(mut),
		untilNotFull:  sync.NewCond(mut),
	}
}

var ErrFull = errors.New("channel is full")
var ErrEmpty = errors.New("channel is empty")

func (ch *BufferedChannel) Send(msg int) {
	ch.mut.Lock()
	defer ch.mut.Unlock() // always defer in case of panic

	for ch.head == ch.tail && ch.isFull {
		ch.untilNotFull.Wait()
		// unlocks the lock put goroutine to sleep until someone wakes it up
		// relocks the lock when you wake up.
	}

	ch.data[ch.head] = msg
	ch.head = (ch.head + 1) % len(ch.data)

	if ch.head == ch.tail {
		ch.isFull = true
	}
	ch.untilNotEmpty.Signal()
}

func (ch *BufferedChannel) Receive() int {
	ch.mut.Lock()
	defer ch.mut.Unlock()

	for ch.head == ch.tail && !ch.isFull {
		ch.untilNotEmpty.Wait()
		// unlocks the lock put goroutine to sleep until someone wakes it up
		// relocks the lock when you wake up.
	}

	msg := ch.data[ch.tail]
	ch.tail = (ch.tail + 1) % len(ch.data)

	ch.isFull = false
	ch.untilNotFull.Signal() // wake up one goroutine waiting on a not-full buffer

	return msg
}
