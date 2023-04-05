package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
	"time"
)

type Barista struct {
	shop *CoffeeShop
	log  *util.Logger
}

func NewBarista(shop *CoffeeShop) *Barista {
	return &Barista{
		shop: shop,
		log:  util.NewLogger("Barista"),
	}
}

// StartWork start the barista working
func (b *Barista) StartWork() {
	go b.doWork()
}

// the barista work loop
func (b *Barista) doWork() {
	loops := 0
	loopTm := time.Now()
	for {
		loops++
		if time.Now().Sub(loopTm) >= 1000*time.Millisecond {
			b.log.Infof("poll %d...", loops)
			loops = 0
			loopTm = time.Now()
		}
		b.CheckCashRegister()
		b.CheckDoneBrewers()
		b.CheckGrinderRefills()
		b.CheckNewOrders()
	}
}

// CheckCashRegister Handle the cash register if a customer is waiting
func (b *Barista) CheckCashRegister() {
	// todo maybe just look at current length of shop's order queue?
	if b.shop.orderObserver.OrdersInThePipe() >= 4 {
		return
	}

	order, ok := b.shop.cashRegister.Barista()
	if !ok {
		return
	}

	beanTypes := model.BeanTypeMap()

	b.log.Infof("got order %v", order)
	order.Start()
	if !beanTypes[order.BeanType] {
		b.log.Infof("bean type unavailable %v", order)
		order.Complete(&model.Receipt{
			Coffee: nil,
			Err:    fmt.Errorf("bean type unavailable %v", order.BeanType),
		})
		return
	}

	b.log.Infof("took order %v", order)

	// seeing an available grinder for an order waiting on the counter is
	// instant... rethink ?
	go func(order *Order) {
		if grinder, ok := <-b.shop.gchan; ok {
			order.SetGrinder(grinder)
			b.shop.orderQueue.Push(order, order.Priority())
		}
	}(order)
}

// CheckDoneBrewers does a non-blocking check for done brewers and put back in available queue
func (b *Barista) CheckDoneBrewers() {
	order, ok := b.shop.brewerDone.Wait0()
	if !ok {
		return
	}
	coffee := order.brewer.GetCoffee()
	b.shop.bchan <- order.brewer // put it back
	order.Complete(&model.Receipt{Coffee: coffee})
}

func (b *Barista) CheckGrinderRefills() {
	select {
	case grinder, ok := <-b.shop.grinderRefill:
		if ok {
			b.log.Infof("grinder refill %v", grinder)
			grinder.Refill(b.shop.roaster)
		}
	default:
	}
}

// CheckNewOrders look for orders that have been paired with a grinder
func (b *Barista) CheckNewOrders() {
	shop := b.shop
	var order *Order

	order, ok := shop.orderQueue.Wait0()
	if !ok {
		return
	}

	extractionProfile := shop.getExtractionProfile(order.StrengthWanted)
	beansNeeded := extractionProfile.GramsFromOunces(order.OuncesOfCoffeeWanted)

	grinder := order.grinder
	groundBeans, err := grinder.Grind(beansNeeded, shop.roaster)
	shop.gchan <- grinder // put it back
	if grinder.ShouldRefill() {
		shop.grinderRefill <- grinder
	}
	if err != nil {
		b.log.Infof("grind error: %v", err) // todo: error handling
		return
	}

	// wait for a brewer
	brewer, ok := <-shop.bchan
	if !ok {
		b.log.Infof("brewers closed") // todo: context shutown system-wide
		return
	}

	order.SetBrewer(brewer)
	brewer.StartBrew(groundBeans, order.OuncesOfCoffeeWanted, func() {
		shop.brewerDone.Push(order, order.Priority())
	})
}
