package queue

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"time"
)

type CashRegister struct {
	pendingOrders chan model.Order
	orderDuration time.Duration
	log           *util.Logger
}

func NewCashRegister(orderTimeMS int) *CashRegister {
	return &CashRegister{
		pendingOrders: make(chan model.Order, 1), // keep this non-0 in size
		orderDuration: time.Duration(orderTimeMS) * time.Millisecond,
		log:           util.NewLogger("CashRegister"),
	}
}

func (c *CashRegister) IsCustomerWaiting() bool {
	return len(c.pendingOrders) != 0
}

func (c *CashRegister) Customer(order model.Order) {
	c.pendingOrders <- order
	c.log.Infof("customer placing order delay %v\n", c.orderDuration)
	time.Sleep(c.orderDuration)
}

func (c *CashRegister) Barista() model.Order {
	order := <-c.pendingOrders
	c.log.Infof("barista taking order delay %v\n", c.orderDuration)
	time.Sleep(c.orderDuration)
	return order
}
