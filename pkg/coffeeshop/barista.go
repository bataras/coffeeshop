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

func (b *Barista) Work() {

	go func() {
		for {
			b.CheckCashRegister()
			b.CheckDoneBrewers()
			b.CheckGrinderRefills()
		}
	}()
}

// CheckCashRegister Handle the cash register if a customer is waiting
func (b *Barista) CheckCashRegister() {
	// todo maybe just look at current length of shop's order queue?
	if b.shop.orderMiddleware.OrdersInThePipe() >= 4 {
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
	b.shop.orderQueue <- order
}

// CheckDoneBrewers non-blocking check for done brewers and put back in available queue
func (b *Barista) CheckDoneBrewers() {
	select {
	case brewer, ok := <-b.shop.brewerDone:
		if ok {
			b.log.Infof("brewer done %v", brewer)
			b.shop.bchan <- brewer // put it back
		}
	default:
	}
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
