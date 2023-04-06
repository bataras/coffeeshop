package coffeeshop

import (
	"sync"
)

// Hopper manages a simple "hopper" as a bounded counter
type Hopper struct {
	mu        sync.Mutex
	beanGrams int
	maxGrams  int
}

func NewHopper(maxGrams int) *Hopper {
	if maxGrams < 0 {
		maxGrams = 0
	}
	return &Hopper{maxGrams: maxGrams}
}

// Count number of grams available
func (h *Hopper) Count() int { return h.beanGrams }

// Size total capacity
func (h *Hopper) Size() int { return h.maxGrams }

// SpaceAvailable how much can we add?
func (h *Hopper) SpaceAvailable() int {
	// don't call internally while holding the mutex
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.maxGrams - h.beanGrams
}

// PercentFull integer percentage 0..100
func (h *Hopper) PercentFull() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Size() == 0 {
		return 0
	}
	return h.Count() * 100 / h.Size()
}

// AddBeans tries to add beans
// returns the actual amount added
func (h *Hopper) AddBeans(grams int) (added int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if grams < 0 {
		grams = 0
	}
	spaceAvail := h.maxGrams - h.beanGrams
	if grams > spaceAvail {
		grams = spaceAvail
	}
	h.beanGrams += grams
	return grams
}

// TakeBeans tries to take the requested amount.
// returns the actual amount taken
func (h *Hopper) TakeBeans(grams int) (got int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if grams < 0 {
		grams = 0
	}
	if grams > h.beanGrams {
		grams = h.beanGrams
	}

	h.beanGrams -= grams
	return grams
}
