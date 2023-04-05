package util

import (
	"container/heap"
	"sync"
)

type PriorityWaitQueue[T any] struct {
	mu       sync.Mutex
	pq       PriorityQueue[T]
	notEmpty chan bool
}

func NewPriorityWaitQueue[T any]() *PriorityWaitQueue[T] {
	ret := &PriorityWaitQueue[T]{
		pq:       make(PriorityQueue[T], 0),
		notEmpty: make(chan bool, 1),
	}
	heap.Init(&ret.pq)
	return ret
}

func (p *PriorityWaitQueue[T]) Push(s T, priority int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pq.PushT(s, priority)
	p.notifyNotEmpty()
}

// Wait waits for an item to be available in the queue
// returns false if the queue is closed
func (p *PriorityWaitQueue[T]) Wait() (T, bool) {
	_, open := <-p.notEmpty
	if !open {
		var zv T
		return zv, false
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	item := p.pq.PopT()
	if p.pq.Len() > 0 {
		p.notifyNotEmpty() // tell the next waiter
	}
	return item, true
}

// Wait0 tries to pop an item from the queue without blocking
// returns false if the queue is empty
func (p *PriorityWaitQueue[T]) Wait0() (T, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.notEmpty) == 0 {
		var zv T
		return zv, false
	}

	<-p.notEmpty

	item := p.pq.PopT()
	if p.pq.Len() > 0 {
		p.notifyNotEmpty() // tell the next waiter
	}
	return item, true
}

// non-blocking signal the notEmpty channel
func (p *PriorityWaitQueue[T]) notifyNotEmpty() {
	select {
	case p.notEmpty <- true:
	default:
	}
}
