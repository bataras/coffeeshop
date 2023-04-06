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

func OrderStrengths() []Strength {
	return []Strength{NormalStrength, MediumStrength, LightStrength}
}

type OrderDoneCB func()

type Order struct {
	OrderNumber          int // Incrementing
	BeanType             string
	OuncesOfCoffeeWanted int
	StrengthWanted       Strength
	done                 chan<- *model.Receipt
	grinder              *Grinder
	brewer               *Brewer
	notifyComplete       func()
}

var orderCount atomic.Int32

func NewOrder(receipts chan<- *model.Receipt, notifyComplete OrderDoneCB) *Order {
	return &Order{
		OrderNumber:    int(orderCount.Add(1)),
		done:           receipts,
		notifyComplete: notifyComplete,
	}
}

func (o *Order) Priority() int {
	return -o.OrderNumber // older orders are higher priority
}

func (o *Order) String() string {
	return fmt.Sprintf("Order#: %d Beans: %v Ounces: %d Strength: %v",
		o.OrderNumber, o.BeanType, o.OuncesOfCoffeeWanted, o.StrengthWanted)
}

func (o *Order) Start() {
}

func (o *Order) Complete(coffee *model.Coffee, err error) {
	o.done <- &model.Receipt{
		OrderNumber: o.OrderNumber,
		Coffee:      coffee,
		Err:         err,
	}
	if o.notifyComplete != nil {
		o.notifyComplete()
	}
}

func (o *Order) SetGrinder(grinder *Grinder) {
	o.grinder = grinder
}

func (o *Order) SetBrewer(brewer *Brewer) {
	o.brewer = brewer
}
