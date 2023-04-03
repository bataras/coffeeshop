package coffeeshop

import (
	"coffeeshop/pkg/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Grinder struct {
	mu                  sync.Mutex
	beanType            BeanType
	hopper              *Hopper
	grindGramsPerSecond util.Rate
	addGramsPerSecond   util.Rate
	refillPercentage    int
}

type BeanGetter func(gramsNeeded int) Beans

func NewGrinder(beanType BeanType, grindGramsPerSecond, addGramsPerSecond, hopperSize int, refillPercentage int) *Grinder {
	val := &Grinder{
		beanType:         beanType,
		hopper:           NewHopper(hopperSize),
		refillPercentage: refillPercentage,
	}
	val.grindGramsPerSecond.SetPerSecond(grindGramsPerSecond)
	val.addGramsPerSecond.SetPerSecond(addGramsPerSecond)
	return val
}

func (g *Grinder) BeanType() BeanType {
	return g.beanType
}

func (g *Grinder) PercentFull() int {
	return g.hopper.PercentFull()
}

// Refill refills the grinder if it's too low on product. takes time.
// try to do this when idle instead of when fulfilling an order
func (g *Grinder) Refill(f BeanGetter) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.refill(f)
}

// do the actual refill. call when locked
func (g *Grinder) refill(f BeanGetter) error {
	if g.hopper.PercentFull() < g.refillPercentage {
		beans := f(g.hopper.SpaceAvailable())
		if beans.beanType != g.BeanType() {
			return fmt.Errorf("tried to refill with wrong beantype")
		}

		g.hopper.AddBeans(beans.weightGrams)

		ms := g.addGramsPerSecond.Duration(beans.weightGrams)
		log.Infof("add beans %v ms %v\n", beans, ms.Milliseconds())
		time.Sleep(ms)
	}
	return nil
}

// Grind grinds beans. takes time
func (g *Grinder) Grind(grams int, f BeanGetter) (Beans, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// ad-hoc refill
	if g.hopper.Count() < grams {
		err := g.refill(f)
		if err != nil {
			return Beans{}, err
		}
	}

	took := g.hopper.TakeBeans(grams)
	if took != grams {
		g.hopper.AddBeans(took)
		return Beans{}, fmt.Errorf("not enough beans. want %v got %v", grams, took)
	}

	ms := g.grindGramsPerSecond.Duration(grams)
	log.Infof("grind beans %v ms %v\n", grams, ms.Milliseconds())
	time.Sleep(ms)
	return Beans{weightGrams: grams}, nil
}
