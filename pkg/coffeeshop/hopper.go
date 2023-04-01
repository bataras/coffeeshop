package coffeeshop

import (
	"fmt"
	"sync"
)

type Hopper struct {
	// We could use a reader/writer lock (RWMutex) but we want this to approximate
	// physical hopper behavior: only one person can take or add beans at a time.
	// todo: This could be enhanced such that beans can "poured in" on top
	// while simultaneously being drained out at the bottom.
	// Doing that should require a Semaphore and making the Add/Take beans methods consume time.
	lock      sync.Mutex
	beanGrams int
}

func NewHopper(startGrams int) Hopper {
	return Hopper{beanGrams: startGrams}
}

func (h *Hopper) AddBeans(grams int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.beanGrams += grams
}

func (h *Hopper) TakeBeans(grams int) (err error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	err = nil
	if h.beanGrams > grams {
		// todo: signal the hopper needs to be filled
		h.beanGrams -= grams
		return nil
	}
	return fmt.Errorf("hopper empty: want %v have %v", grams, h.beanGrams)
}
