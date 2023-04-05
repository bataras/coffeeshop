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

func (o *OrderObserver) OrderTaken(order *Order) {
	o.ordersInThePipe.Add(1)
	o.log.Infof("order taken %v", order)
}

func (o *OrderObserver) OrderCompleted(order *Order) {
	o.ordersInThePipe.Add(-1)
	o.log.Infof("order complete %v", order)
}

func (o *OrderObserver) UsingGrinder(order *Order) {
	o.log.Infof("order has grinder %v", order)
}

func (o *OrderObserver) UsingBrewer(order *Order) {
	o.log.Infof("order has brewer %v", order)
}
