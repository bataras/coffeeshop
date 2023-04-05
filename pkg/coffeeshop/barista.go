package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
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
	for {
		b.log.Infof("wait for work...")

		var order *Order
		var grinder *Grinder
		ok := false

		select {
		case order, ok = <-b.shop.cashRegister.GetWaitChan():
			if ok {
				b.HandleOrderFromCashRegister(order)
			}
		case _, ok = <-b.shop.orderQueue.GetWaitChan():
			if ok {
				order, ok = b.shop.orderQueue.Pop()
				if ok {
					b.HandleNewOrder(order)
				}
			}
		case _, ok = <-b.shop.brewerDone.GetWaitChan():
			if ok {
				order, ok = b.shop.brewerDone.Pop()
				if ok {
					b.HandleDoneBrewer(order)
				}
			}
		case grinder, ok = <-b.shop.grinderRefill:
			if ok {
				b.HandleGrinderRefill(grinder)
			}
		}
	}
}

// HandleOrderFromCashRegister Handle the cash register if a customer is waiting
func (b *Barista) HandleOrderFromCashRegister(order *Order) {
	// todo maybe just look at current length of shop's order queue?
	// if b.shop.orderObserver.OrdersInThePipe() >= 4 {
	// 	return
	// }

	// barista is doing work here
	b.shop.cashRegister.SpendTimeHandlingAnOrder(false)

	beanTypes := model.BeanTypeMap()

	b.log.Infof("took order %v", order)
	order.Start()
	if !beanTypes[order.BeanType] {
		b.log.Infof("bean type unavailable %v", order)
		order.Complete(nil, fmt.Errorf("bean type unavailable %v", order.BeanType))
		return
	}

	// seeing an available grinder for an order waiting on the counter is
	// instant... rethink ?
	go func(order *Order) {
		if grinder, ok := <-b.shop.gchan; ok {
			order.SetGrinder(grinder)
			b.shop.orderQueue.Push(order, order.Priority())
		} else {
			order.Complete(nil, fmt.Errorf("grinders are closed"))
		}
	}(order)
}

// CheckDoneBrewers does a non-blocking check for done brewers and put back in available queue
func (b *Barista) HandleDoneBrewer(order *Order) {
	b.log.Infof("brewer done %v", order)
	coffee := order.brewer.GetCoffee()
	b.shop.bchan <- order.brewer // put it back
	order.Complete(coffee, nil)
}

func (b *Barista) HandleGrinderRefill(grinder *Grinder) {
	b.log.Infof("grinder refill %v", grinder)
	grinder.Refill(b.shop.roaster)
}

// HandleNewOrders handle orders that have been paired with a grinder
func (b *Barista) HandleNewOrder(order *Order) {
	shop := b.shop

	b.log.Infof("grind start %v", order)

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
		order.Complete(nil, err)
		return
	}

	// wait for a brewer
	brewer, ok := <-shop.bchan
	if !ok {
		b.log.Infof("brewers closed") // todo: context shutown system-wide
		order.Complete(nil, fmt.Errorf("brewers are closed"))
		return
	}

	order.SetBrewer(brewer)
	brewer.StartBrew(groundBeans, order.OuncesOfCoffeeWanted, func() {
		shop.brewerDone.Push(order, order.Priority())
	})
}
