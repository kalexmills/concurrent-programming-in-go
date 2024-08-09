package main

import (
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
			bc.Send(i)
			fmt.Printf("sent %d\n", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			d := bc.Receive()
			fmt.Printf("received %d\n", d)
		}
	}()

	wg.Wait()
}

type BufferedChannel struct {
	data   []int
	head   int // always pointing to next free space
	tail   int // always pointed to next item to be received.
	isFull bool

	mut      *sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

//   h
// [ 3 4 5 6 ]
//   t

func NewBufferedChannel(size int) *BufferedChannel {
	mut := &sync.Mutex{}
	return &BufferedChannel{
		data:     make([]int, size),
		mut:      mut,
		notFull:  sync.NewCond(mut),
		notEmpty: sync.NewCond(mut),
	}
}

func (bc *BufferedChannel) Send(data int) {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.isFull {
		bc.notFull.Wait() // unlocks bc.mut
	}

	bc.data[bc.head] = data
	bc.head = (bc.head + 1) % len(bc.data)

	if bc.head == bc.tail {
		bc.isFull = true
	}

	bc.notEmpty.Signal()
}

func (bc *BufferedChannel) Receive() int {
	bc.mut.Lock()
	defer bc.mut.Unlock()

	for bc.tail == bc.head && !bc.isFull {
		bc.notEmpty.Wait()
	}

	result := bc.data[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.data)

	bc.isFull = false
	bc.notFull.Signal()

	return result
}
