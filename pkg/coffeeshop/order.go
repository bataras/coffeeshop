package coffeeshop

import "coffeeshop/pkg/model"

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
	BeanType             model.BeanType
	OuncesOfCoffeeWanted int
	StrengthWanted       Strength
	done                 chan<- *model.Receipt
	observer             IOrderObserver
	grinder              *Grinder
	brewer               *Brewer
	// todo: maybe have a audit/observable mechanism and return the order to the customer instead of the receipt channel
}

func NewOrder(receipts chan<- *model.Receipt, orderMiddleware IOrderObserver) *Order {
	return &Order{
		done:     receipts,
		observer: orderMiddleware,
	}
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
