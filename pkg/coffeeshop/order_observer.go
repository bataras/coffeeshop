package coffeeshop

import (
	"coffeeshop/pkg/util"
	"sync/atomic"
)

var _ IOrderObserver = &OrderObserver{}

type IOrderObserver interface {
	OrderTaken(*Order)
	OrderCompleted(*Order)
	UsingGrinder(*Order)
	UsingBrewer(*Order)
	OrdersInThePipe() int
}

type OrderObserver struct {
	log             *util.Logger
	ordersInThePipe atomic.Int32
}

func NewOrderObserver() IOrderObserver {
	return &OrderObserver{
		log: util.NewLogger("OrderObserver"),
	}
}

func (o *OrderObserver) OrdersInThePipe() int {
	return int(o.ordersInThePipe.Load())
}

func (o *OrderObserver) OrderTaken(_ *Order) {
	o.ordersInThePipe.Add(1)
	o.log.Infof("order taken. pipe=%v", o.OrdersInThePipe())
}

func (o *OrderObserver) OrderCompleted(_ *Order) {
	o.ordersInThePipe.Add(-1)
	o.log.Infof("order complete. pipe=%v", o.OrdersInThePipe())
}

func (o *OrderObserver) UsingGrinder(_ *Order) {
	o.log.Infof("order has grinder. pipe=%v", o.OrdersInThePipe())
}

func (o *OrderObserver) UsingBrewer(_ *Order) {
	o.log.Infof("order has brewer. pipe=%v", o.OrdersInThePipe())
}
