package coffeeshop

import (
	"coffeeshop/pkg/util"
	"sync/atomic"
	"time"
)

type CashRegister struct {
	pendingOrders chan *Order
	baristaCount  atomic.Int32
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

// Barista pretends to be a barista tending to the cash register
// will give up after a timeout if no customers are waiting
// only 1 barista at a time can be at the register
func (c *CashRegister) Barista(timeoutMs int) (*Order, bool) {
	// only 1 barista at a time can be at the register
	if c.baristaCount.Add(1) > 2 {
		c.baristaCount.Add(-1)
		return nil, true
	}
	defer c.baristaCount.Add(-1)

	select {
	case order, ok := <-c.pendingOrders:
		return order, ok
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		c.log.Infof("barista timeout")
	}

	return nil, true
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
