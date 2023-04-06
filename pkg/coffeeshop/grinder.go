package coffeeshop

import (
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Grinder struct {
	mu                  sync.Mutex
	log                 *util.Logger
	id                  int
	beanType            string
	hopper              *Hopper
	grindGramsPerSecond util.Rate
	addGramsPerSecond   util.Rate
	refillGate          atomic.Int32
	refillPercentage    int
}

type IRoaster interface {
	GetBeans(gramsNeeded int, beanType string) model.Beans
}

type IRoasterFunc func(gramsNeeded int, beanType string) model.Beans

func (f IRoasterFunc) GetBeans(gramsNeeded int, beanType string) model.Beans {
	return f(gramsNeeded, beanType)
}

var grinderCount atomic.Int32

func NewGrinder(cfg *config.GrinderCfg) *Grinder {
	num := int(orderCount.Add(1))
	val := &Grinder{
		log:              util.NewLogger(fmt.Sprintf("Grinder %d %s", num, cfg.BeanCfg.BeanType)),
		id:               num,
		beanType:         cfg.BeanCfg.BeanType,
		hopper:           NewHopper(cfg.HopperSize),
		refillPercentage: cfg.RefillPercentage,
	}
	val.grindGramsPerSecond.SetPerSecond(cfg.GrindGramsPerSecond)
	val.addGramsPerSecond.SetPerSecond(cfg.AddGramsPerSecond)
	return val
}

func (g *Grinder) String() string {
	return fmt.Sprintf("bean: %v hopper %v %v", g.beanType, g.hopper.Count(), g.hopper.PercentFull())
}

func (g *Grinder) BeanType() string {
	return g.beanType
}

func (g *Grinder) ShouldRefill() bool {
	return g.hopper.PercentFull() < g.refillPercentage
}

// TryRefill refills the grinder if it's too low on product. takes time.
// try to do this when idle instead of when fulfilling an order so that
// beans are ready when the customer orders
func (g *Grinder) TryRefill(f IRoaster) error {
	g.log.Infof("try background refill")

	// This method is called externally to try to refill the hopper during idle time
	// Use this simple one-way gate to avoid waiting for a refill that may already
	// be happening on another thread ad-hoc while grinding
	g.refillGate.Add(-1)
	if g.refillGate.Add(1) > 0 {
		return nil
	}

	// if we sneak through here, another thread that wants to grind will have to wait
	// for this background refill. Or in a rare case, this thread will wait for the
	// other thread only to find the refill is done.
	// That's ok. It's kind of like a mutex with a weak lock path and a strong lock path.
	// The weak lock (this one) may abort without waiting, but the strong lock (Grind)
	// will always wait to succeed.
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.refillInternal(f, false)
}

// do the actual refill. call when locked
func (g *Grinder) refillInternal(roaster IRoaster, adHoc bool) error {
	if g.hopper.PercentFull() >= g.refillPercentage {
		g.log.Infof("refill not necessary: adhoc %v", adHoc)
		return nil
	}

	g.log.Infof("refill: adhoc %v", adHoc)

	beans := roaster.GetBeans(g.hopper.SpaceAvailable(), g.beanType)
	if beans.BeanType != g.beanType {
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
	g.refillGate.Add(1)
	defer g.refillGate.Add(-1)

	// by coffee shop design, this method is never supposed to be reentered. But we do
	// want to prevent hopper refills from being reentered. It's just safer to
	// lock the whole Grind method rather than just the internal refill
	g.mu.Lock()
	defer g.mu.Unlock()

	// ad-hoc refill
	if g.hopper.Count() < grams {
		err := g.refillInternal(roaster, true)
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
	return model.Beans{BeanType: g.beanType, WeightGrams: grams}, nil
}
