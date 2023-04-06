package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
	"sync/atomic"
	"time"
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
	notifyComplete       OrderDoneCB
	started              time.Time
	log                  *util.Logger
}

var orderCount atomic.Int32

func NewOrder(receipts chan<- *model.Receipt, notifyComplete OrderDoneCB) *Order {
	num := int(orderCount.Add(1))
	return &Order{
		OrderNumber:    num,
		done:           receipts,
		notifyComplete: notifyComplete,
		started:        time.Now(),
		log:            util.NewLogger(fmt.Sprintf("--- Order %d", num)),
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
	o.started = time.Now()
	o.log.Infof("start")
}

func (o *Order) Complete(coffee *model.Coffee, err error) {
	took := time.Now().Sub(o.started)
	if err == nil {
		o.log.Infof("complete. took %v", took)
	} else {
		o.log.Errorf("complete. took %v: %v", took, err)
	}
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
	o.log.Infof("set grinder: %v", grinder)
	o.grinder = grinder
}

func (o *Order) SetBrewer(brewer *Brewer) {
	o.log.Infof("set brewer")
	o.brewer = brewer
}
