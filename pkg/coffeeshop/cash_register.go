package coffeeshop

import (
	"coffeeshop/pkg/util"
	"time"
)

type CashRegister struct {
	pendingOrders chan *Order
	orderDuration time.Duration
	log           *util.Logger
}

func NewCashRegister(orderTimeMS int) *CashRegister {
	return &CashRegister{
		pendingOrders: make(chan *Order), // should be non-buffered
		orderDuration: time.Duration(orderTimeMS) * time.Millisecond,
		log:           util.NewLogger("CashRegister"),
	}
}

// Customer blocks, waiting for a barista
func (c *CashRegister) Customer(order *Order) {
	c.pendingOrders <- order
	c.SpendTimeHandlingAnOrder(true)
}

func (c *CashRegister) SpendTimeHandlingAnOrder(asCustomer bool) {
	if asCustomer {
		c.log.Infof("customer placing order delay %v", c.orderDuration)
	} else {
		c.log.Infof("barista taking order delay %v", c.orderDuration)
	}
	time.Sleep(c.orderDuration)
}

func (c *CashRegister) GetWaitChan() <-chan *Order {
	return c.pendingOrders
}
