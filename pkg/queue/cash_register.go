package queue

import (
	"coffeeshop/pkg/util"
	"time"
)

type CashRegister struct {
	orderTaker chan bool
	orderTime  time.Duration
	log        *util.Logger
}

func NewCashRegister(orderTimeMS int) *CashRegister {
	return &CashRegister{
		orderTaker: make(chan bool),
		orderTime:  time.Duration(orderTimeMS) * time.Millisecond,
		log:        util.NewLogger("CashRegister"),
	}
}

func (c *CashRegister) Customer() {
	c.orderTaker <- true
	c.log.Infof("customer placing order delay %v\n", c.orderTime)
	time.Sleep(c.orderTime)
}

func (c *CashRegister) Barista() {
	<-c.orderTaker
	c.log.Infof("barista taking order delay %v\n", c.orderTime)
	time.Sleep(c.orderTime)
}
