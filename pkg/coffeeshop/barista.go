package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
	"sync/atomic"
)

type Barista struct {
	id   int
	shop *CoffeeShop
	log  *util.Logger
}

var baristaCount atomic.Int32

func NewBarista(shop *CoffeeShop) *Barista {
	id := int(baristaCount.Add(1))
	return &Barista{
		id:   id,
		shop: shop,
		log:  util.NewLogger(fmt.Sprintf("Barista %v", id)),
	}
}

// StartWork start the barista working
func (b *Barista) StartWork() {
	go b.doWork()
}

// the barista work loop
func (b *Barista) doWork() {
	for {
		// b.log.Infof("wait for work...")

		var order *Order
		var grinder *Grinder
		ok := false
		select {
		case order, ok = <-b.shop.cashRegister.GetWaitChan():
			if ok {
				b.log.Infof("handle cash register with %d orders in the pipe", len(b.shop.orderPipeDepth))
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
				b.shop.grinders <- grinder // put it back in rotation
			}
		}
	}
}

// HandleOrderFromCashRegister Handle the cash register if a customer is waiting
func (b *Barista) HandleOrderFromCashRegister(order *Order) {
	// barista is doing work here, talking to the customer
	b.shop.cashRegister.SpendTimeHandlingAnOrder(false)

	beanTypes := model.BeanTypeMap()

	b.log.Infof("took order %v", order)
	order.Start()
	if !beanTypes[order.BeanType] {
		b.log.Infof("bean type unavailable %v", order)
		order.Complete(nil, fmt.Errorf("bean type unavailable %v", order.BeanType))
		return
	}

	// seeing an available grinder for an order waiting on the counter is essentially a signal
	go func(order *Order) {
		if grinder, ok := <-b.shop.grinders; ok {
			order.SetGrinder(grinder)
			b.shop.orderQueue.Push(order, order.Priority())
		} else {
			order.Complete(nil, fmt.Errorf("grinders are closed"))
		}
	}(order)
}

// HandleNewOrder handle orders that have been paired with a grinder
func (b *Barista) HandleNewOrder(order *Order) {
	shop := b.shop

	b.log.Infof("grind start %v", order)

	extractionProfile := shop.getExtractionProfile(order.StrengthWanted)
	beansNeeded := extractionProfile.GramsFromOunces(order.OuncesOfCoffeeWanted)

	grinder := order.grinder
	groundBeans, err := grinder.Grind(beansNeeded, shop.roaster)
	if grinder.ShouldRefill() {
		shop.grinderRefill <- grinder
	} else {
		shop.grinders <- grinder // put it back in rotation
	}
	if err != nil {
		b.log.Infof("grind error: %v", err)
		order.Complete(nil, err)
		return
	}

	// seeing an available brewer for an order waiting on the counter is essentially a signal
	go func(order *Order) {
		if brewer, ok := <-b.shop.brewers; ok {
			order.SetBrewer(brewer)
			brewer.StartBrew(groundBeans, order.OuncesOfCoffeeWanted, func() {
				shop.brewerDone.Push(order, order.Priority())
			})
		} else {
			b.log.Infof("brewers closed")
			order.Complete(nil, fmt.Errorf("brewers are closed"))
		}
	}(order)
}

// HandleDoneBrewer does a non-blocking check for done brewers and put back in available queue
func (b *Barista) HandleDoneBrewer(order *Order) {
	b.log.Infof("brewer done %v", order)
	coffee := order.brewer.GetCoffee()
	b.shop.brewers <- order.brewer // put it back
	order.Complete(coffee, nil)
}

func (b *Barista) HandleGrinderRefill(grinder *Grinder) {
	b.log.Infof("grinder refill %v", grinder)
	if err := grinder.Refill(b.shop.roaster); err != nil {
		b.log.Infof("grinder refill error %v", err)
	}
}
