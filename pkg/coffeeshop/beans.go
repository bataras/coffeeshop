package coffeeshop

type BeanType int

const (
	Columbian BeanType = iota
	Ethiopian
	French
	Italian
)

func (b BeanType) String() string {
	switch b {
	case Columbian:
		return "Columbian"
	case Ethiopian:
		return "Ethiopian"
	case French:
		return "French"
	case Italian:
		return "Italian"
	default:
		return "Unknown"
	}
}

type Beans struct {
	beanType    BeanType
	weightGrams int
	// indicate some state change? create a new type?
}
