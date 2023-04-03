package coffeeshop

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
	OuncesOfCoffeeWanted int
	StrengthWanted       Strength
}

type Receipt struct {
	Coffee *Coffee
	Err    error
}
