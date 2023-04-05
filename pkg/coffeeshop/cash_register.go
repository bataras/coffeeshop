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
	c.log.Infof("customer placing order delay %v\n", c.orderDuration)
	time.Sleep(c.orderDuration)
}

// Barista doesn't block if there are no orders waiting
func (c *CashRegister) Barista() (*Order, bool) {
	// todo: only allow 1 barista per register
	select {
	case order, ok := <-c.pendingOrders:
		if !ok {
			return nil, ok
		}
		c.log.Infof("barista taking order delay %v\n", c.orderDuration)
		time.Sleep(c.orderDuration)
		return order, ok

	default:
		return nil, false
	}
}
