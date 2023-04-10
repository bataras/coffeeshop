package util

import (
	"sync"
)

var _ ITryMutex = &TryLock{}
var _ ITryMutex = &TryLockC{}

type ITryMutex interface {
	TryLock() bool
	Lock()
	Unlock()
}

/*
TryLockC implements a try-able mutex using a channel
*/
type TryLockC struct {
	lock chan bool
}

/*
TryLock implements a try-able mutex a counter gate
*/
type TryLock struct {
	gate  sync.Mutex
	mu    sync.Mutex
	count int
}

/* -------------------- TryLockC -------------------- */

func (t *TryLock) TryLock() bool {
	t.gate.Lock()
	if t.count != 0 { // safe to lock?
		// fmt.Printf("can't lock. count %v\n", t.count)
		t.gate.Unlock()
		return false
	}
	t.count++

	t.mu.Lock()
	t.gate.Unlock()
	// fmt.Printf("trylock count %v %v\n", t.count, t.mu)
	return true
}

func (t *TryLock) Lock() {
	t.gate.Lock()
	t.count++
	t.gate.Unlock()
	t.mu.Lock()
}

func (t *TryLock) Unlock() {
	// fmt.Printf("unlock %v %v\n", t.count, t.mu)
	t.count--
	t.mu.Unlock()
}

/* -------------------- TryLockC -------------------- */

func NewTryLockC() *TryLockC {
	return &TryLockC{lock: make(chan bool, 1)}
}

func (t *TryLockC) TryLock() bool {
	select {
	case t.lock <- true:
	default:
		return false
	}
	return true
}

func (t *TryLockC) Lock() {
	t.lock <- true
}

func (t *TryLockC) Unlock() {
	<-t.lock
}
