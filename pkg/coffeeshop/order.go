package coffeeshop

import (
	"coffeeshop/pkg/model"
	"fmt"
	"sync/atomic"
)

type Strength int

const (
	NormalStrength Strength = iota
	MediumStrength
	LightStrength
)

func (s Strength) String() string {
	switch s {
	case NormalStrength:
		return "NormalStrength"
	case MediumStrength:
		return "MediumStrength"
	case LightStrength:
		return "LightStrength"
	default:
		return "Unknown Strength"
	}
}

type Order struct {
	OrderNumber          int // Incrementing
	BeanType             model.BeanType
	OuncesOfCoffeeWanted int
	StrengthWanted       Strength
	done                 chan<- *model.Receipt
	observer             IOrderObserver
	grinder              *Grinder
	brewer               *Brewer
	// todo: maybe have a audit/observable mechanism and return the order to the customer instead of the receipt channel
}

var orderCount atomic.Int32

func NewOrder(receipts chan<- *model.Receipt, orderMiddleware IOrderObserver) *Order {
	return &Order{
		OrderNumber: int(orderCount.Add(1)),
		done:        receipts,
		observer:    orderMiddleware,
	}
}

func (o *Order) Priority() int {
	return -o.OrderNumber // older orders are higher priority
}

func (o *Order) String() string {
	return fmt.Sprintf("No: %d Beans: %v Ounces: %d Strength: %v",
		o.OrderNumber, o.BeanType, o.OuncesOfCoffeeWanted, o.StrengthWanted)
}

func (o *Order) Start() {
	o.observer.OrderTaken(o)
}

func (o *Order) Complete(receipt *model.Receipt) {
	o.observer.OrderCompleted(o)
	o.done <- receipt
}

func (o *Order) SetGrinder(grinder *Grinder) {
	o.grinder = grinder
	o.observer.UsingGrinder(o)
}

func (o *Order) SetBrewer(brewer *Brewer) {
	o.brewer = brewer
	o.observer.UsingBrewer(o)
}
