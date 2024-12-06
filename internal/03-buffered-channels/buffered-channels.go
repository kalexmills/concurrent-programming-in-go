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
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			err := ch.Send(i)
			fmt.Printf("sent %d, err = %v\n", i, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			data, err := ch.Receive()
			fmt.Printf("received %d, err = %v\n", data, err)
		}
	}()

	wg.Wait()
}

type BufferedChan struct {
	data   []int
	head   int // next empty space
	tail   int // next space to read
	isFull bool

	mut *sync.Mutex
}

func NewBufferedChan(size int) *BufferedChan {
	return &BufferedChan{
		data: make([]int, size),
		mut:  &sync.Mutex{},
	}
}

var ErrFull = errors.New("full")
var ErrEmpty = errors.New("empty")

func (bc *BufferedChan) Send(msg int) error {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	if bc.head == bc.tail && bc.isFull {
		return ErrFull
	}

	bc.data[bc.head] = msg
	bc.head = (bc.head + 1) % len(bc.data)

	if bc.head == bc.tail {
		bc.isFull = true
	}

	return nil
}

func (bc *BufferedChan) Receive() (int, error) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	if bc.head == bc.tail && !bc.isFull {
		return 0, ErrEmpty
	}

	msg := bc.data[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.data)

	bc.isFull = false

	return msg, nil
}
