package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
	"sync"
	"time"
)

type Grinder struct {
	mu                  sync.Mutex
	log                 *util.Logger
	beanType            model.BeanType
	hopper              *Hopper
	grindGramsPerSecond util.Rate
	addGramsPerSecond   util.Rate
	refillPercentage    int
}

type IRoaster interface {
	GetBeans(gramsNeeded int, beanType model.BeanType) model.Beans
}

type IRoasterFunc func(gramsNeeded int) model.Beans

func (f IRoasterFunc) GetBeans(gramsNeeded int) model.Beans {
	return f(gramsNeeded)
}

func NewGrinder(beanType model.BeanType, grindGramsPerSecond,
	addGramsPerSecond, hopperSize int, refillPercentage int) *Grinder {
	val := &Grinder{
		log:              util.NewLogger("Grinder"),
		beanType:         beanType,
		hopper:           NewHopper(hopperSize),
		refillPercentage: refillPercentage,
	}
	val.grindGramsPerSecond.SetPerSecond(grindGramsPerSecond)
	val.addGramsPerSecond.SetPerSecond(addGramsPerSecond)
	return val
}

func (g *Grinder) String() string {
	return fmt.Sprintf("bean: %v hopper %v %v", g.beanType, g.hopper.Count(), g.hopper.PercentFull())
}

func (g *Grinder) BeanType() model.BeanType {
	return g.beanType
}

func (g *Grinder) ShouldRefill() bool {
	return g.hopper.PercentFull() < g.refillPercentage
}

// Refill refills the grinder if it's too low on product. takes time.
// try to do this when idle instead of when fulfilling an order
func (g *Grinder) Refill(f IRoaster) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.refillInternal(f)
}

// do the actual refill. call when locked
func (g *Grinder) refillInternal(roaster IRoaster) error {
	if g.hopper.PercentFull() >= g.refillPercentage {
		g.log.Infof("refill not necessary")
		return nil
	}

	beans := roaster.GetBeans(g.hopper.SpaceAvailable(), g.beanType)
	if beans.BeanType != g.BeanType() {
		return fmt.Errorf("tried to refill with wrong beantype")
	}

	g.hopper.AddBeans(beans.WeightGrams)

	ms := g.addGramsPerSecond.Duration(beans.WeightGrams)
	g.log.Infof("add beans %v ms %v", beans, ms.Milliseconds())
	time.Sleep(ms)
	return nil
}

// Grind grinds beans. takes time
func (g *Grinder) Grind(grams int, roaster IRoaster) (model.Beans, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// ad-hoc refill
	if g.hopper.Count() < grams {
		err := g.refillInternal(roaster)
		if err != nil {
			return model.Beans{}, err
		}
	}

	took := g.hopper.TakeBeans(grams)
	if took != grams {
		g.hopper.AddBeans(took)
		// requested grams will never be satisfied...
		return model.Beans{}, fmt.Errorf("not enough beans. want %v got %v", grams, took)
	}

	ms := g.grindGramsPerSecond.Duration(grams)
	g.log.Infof("grind beans %v ms %v", grams, ms.Milliseconds())
	time.Sleep(ms)
	return model.Beans{WeightGrams: grams}, nil
}
