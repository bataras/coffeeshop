package model

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
	BeanType             BeanType
	OuncesOfCoffeeWanted int
	StrengthWanted       Strength
	done                 chan<- *Receipt
	// todo: maybe have a audit/observable mechanism and return the order to the customer instead of the receipt channel
}

func NewOrder(receipts chan<- *Receipt) Order {
	return Order{done: receipts}
}

func (o *Order) Complete(receipt *Receipt) {
	o.done <- receipt
}

type Receipt struct {
	Coffee *Coffee
	Err    error
}
