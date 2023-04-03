package coffeeshop

import "fmt"

type Coffee struct {
	beanType BeanType
	ounces   int
}

func NewCoffee(beanType BeanType, ounces int) *Coffee {
	return &Coffee{
		beanType: beanType,
		ounces:   ounces,
	}
}

func (c *Coffee) String() string {
	return fmt.Sprintf("{%v %v}", c.beanType, c.ounces)
}

func (c *Coffee) BeanType() BeanType {
	return c.beanType
}

func (c *Coffee) Ounces() int {
	return c.ounces
}
