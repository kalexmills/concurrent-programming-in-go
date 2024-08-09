package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	bc := NewBufferedChannel(4)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			err := bc.Send(i)
			fmt.Printf("sent %d with err: %v\n", i, err)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			d, err := bc.Receive()
			fmt.Printf("received %d with err: %v\n", d, err)
		}
	}()

	wg.Wait()
}

type BufferedChannel struct {
	data   []int
	head   int // always pointing to next free space
	tail   int // always pointed to next item to be received.
	isFull bool

	mut *sync.Mutex
}

//   h
// [ 3 4 5 6 ]
//   t

func NewBufferedChannel(size int) *BufferedChannel {
	return &BufferedChannel{
		data: make([]int, size),
		mut:  &sync.Mutex{},
	}
}

var ErrFull = errors.New("buffer full")
var ErrEmpty = errors.New("buffer empty")

func (bc *BufferedChannel) Send(data int) error {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	if bc.isFull {
		return ErrFull
	}

	bc.data[bc.head] = data
	bc.head = (bc.head + 1) % len(bc.data)

	if bc.head == bc.tail {
		bc.isFull = true
	}

	return nil
}

func (bc *BufferedChannel) Receive() (int, error) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	if bc.tail == bc.head && !bc.isFull {
		return 0, ErrEmpty
	}

	result := bc.data[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.data)

	bc.isFull = false

	return result, nil
}
