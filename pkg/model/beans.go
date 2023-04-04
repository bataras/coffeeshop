package model

type BeanType int

const (
	Columbian BeanType = iota
	Ethiopian
	French
	Italian
)

func BeanTypeList() []BeanType {
	return []BeanType{Columbian, Ethiopian, French, Italian}
}

func BeanTypeMap() (ret map[BeanType]bool) {
	ret = map[BeanType]bool{}
	for _, t := range BeanTypeList() {
		ret[t] = true
	}
	return
}

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
	BeanType    BeanType
	WeightGrams int
	// indicate some state change? create a new type?
}
