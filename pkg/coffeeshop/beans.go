package coffeeshop

type BeanType int

const (
	Columbian BeanType = iota
	Ethiopian
	French
	Italian
)

type Beans struct {
	beanType    BeanType
	weightGrams int
	// indicate some state change? create a new type?
}
