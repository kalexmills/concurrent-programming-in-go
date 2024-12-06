package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	ch := NewBufferedChan(4) // make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	for j := 0; j < 2; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				ch.Send(i)
				fmt.Printf("sent %d\n", i)
			}
		}()
	}

	for j := 0; j < 2; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				data := ch.Receive()
				fmt.Printf("received %d\n", data)
			}
		}()
	}

	wg.Wait()
}

type BufferedChan struct {
	data   []int
	head   int // next empty space
	tail   int // next space to read
	isFull bool

	mut             *sync.Mutex
	waitForNotFull  *sync.Cond
	waitForNotEmpty *sync.Cond
}

func NewBufferedChan(size int) *BufferedChan {
	mut := &sync.Mutex{}
	return &BufferedChan{
		data:            make([]int, size),
		mut:             mut,
		waitForNotEmpty: sync.NewCond(mut),
		waitForNotFull:  sync.NewCond(mut),
	}
}

var ErrFull = errors.New("full")
var ErrEmpty = errors.New("empty")

func (bc *BufferedChan) Send(msg int) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && bc.isFull {
		bc.waitForNotFull.Wait() // unlock; wait; then get lock back
	}
	// we have the lock; and buffer is not full

	bc.data[bc.head] = msg
	bc.head = (bc.head + 1) % len(bc.data)

	if bc.head == bc.tail {
		bc.isFull = true
	}
	bc.waitForNotEmpty.Signal()
}

func (bc *BufferedChan) Receive() int {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.head == bc.tail && !bc.isFull {
		bc.waitForNotEmpty.Wait()
	}

	msg := bc.data[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.data)

	bc.isFull = false
	bc.waitForNotFull.Signal() // wake up exactly one goroutine to make progress

	return msg
}
