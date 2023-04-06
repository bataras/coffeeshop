package coffeeshop

import (
	"fmt"
	"strings"
	"sync"
)

type GrinderPool struct {
	mu          sync.Mutex
	grinders    map[string]chan *Grinder // map by bean type
	maxGrinders int
}

func NewGrinderPool(maxGrinders int) *GrinderPool {
	return &GrinderPool{
		grinders:    map[string]chan *Grinder{},
		maxGrinders: maxGrinders,
	}
}

func (p *GrinderPool) Put(grinder *Grinder) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bean := strings.ToLower(grinder.BeanType())
	if _, have := p.grinders[bean]; !have {
		p.grinders[bean] = make(chan *Grinder, p.maxGrinders)
	}

	if ch, have := p.grinders[bean]; have { // sanity
		ch <- grinder
	}
}

// ChanFor waits for a grinder to be available
func (p *GrinderPool) ChanFor(beanType string) (<-chan *Grinder, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bean := strings.ToLower(beanType)
	if ch, have := p.grinders[bean]; !have {
		return nil, fmt.Errorf("no grinder pool for: %v", bean)
	} else {
		return ch, nil
	}
}
