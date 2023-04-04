package middleware

import (
	"coffeeshop/pkg/model"
	"sync/atomic"
)

type Orders struct {
	ordersInThePipe atomic.Int32
}

func NewOrders() *Orders {
	return &Orders{}
}

func (o *Orders) OrdersInThePipe() int {
	return int(o.ordersInThePipe.Load())
}

func (o *Orders) OrderTaken(_ *model.Order) {
	o.ordersInThePipe.Add(1)
}

func (o *Orders) OrderCompleted(_ *model.Order) {
	o.ordersInThePipe.Add(-1)
}
