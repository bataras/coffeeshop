package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewPriorityWaitQueue(t *testing.T) {
	q := NewPriorityWaitQueue[string]()

	_, notEmpty := q.Wait0()
	assert.False(t, notEmpty)

	q.Push("aaa", 1)

	aaa, notEmpty2 := q.Wait0()
	assert.True(t, notEmpty2)
	assert.Equal(t, "aaa", aaa)

	q.Push("a", 1)
	q.Push("b", 2)
	q.Push("c", 3)
	q.Push("d", 4)
	q.Push("e", 5)

	wg := sync.WaitGroup{}
	fire := make(chan bool)
	listener := func() {
		<-fire
		str, ok := q.Wait()
		fmt.Printf("got %v %v\n", str, ok)
		wg.Done()
	}
	wg.Add(3)
	go listener()
	go listener()
	go listener()

	close(fire)
	wg.Wait()

	str, _ := q.Wait()
	assert.Equal(t, "b", str)
	str, _ = q.Wait()
	assert.Equal(t, "a", str)
}

func TestNewPriorityWaitQueue2(t *testing.T) {
	q := NewPriorityWaitQueue[string]()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, notEmpty := q.Wait0()
		fmt.Printf("waiting %v...\n", notEmpty)
		wg.Done()
		str, ok := q.Wait()
		fmt.Printf("got %v %v\n", str, ok)
		wg.Done()
	}()

	wg.Wait()
	fmt.Printf("started...\n")
	wg.Add(1)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("push blah...\n")
	q.Push("blah", 4)
	wg.Wait()
	fmt.Printf("received...\n")
}
