package util

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTryLock_Lock(t1 *testing.T) {
	hammerITryMutex_Lock(t1, &TryLock{})
	hammerITryMutex_Lock(t1, NewTryLockC())
}

func hammerITryMutex_Lock(t1 *testing.T, lock ITryMutex) {
	fire := make(chan bool)
	wg := sync.WaitGroup{}

	trier := func() {
		<-fire
		time.Sleep(100 * time.Millisecond)
		fmt.Println("trier try")
		tm := time.Now()
		cnt := 0
		for {
			cnt++
			if lock.TryLock() {
				fmt.Printf("try lock after %v tries in %v\n", cnt, time.Now().Sub(tm))
				lock.Unlock()
				break
			}
		}
		fmt.Println("trier out")
		wg.Done()
	}

	locker := func() {
		<-fire
		lock.Lock()
		fmt.Println("lock in")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("lock out")
		lock.Unlock()
		wg.Done()
	}

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go trier()
		go locker()
	}

	close(fire)
	wg.Wait()
}
