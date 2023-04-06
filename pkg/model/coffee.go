package model

import "fmt"

type Coffee struct {
	beanType string
	ounces   int
}

func NewCoffee(beanType string, ounces int) *Coffee {
	return &Coffee{
		beanType: beanType,
		ounces:   ounces,
	}
}

func (c *Coffee) String() string {
	return fmt.Sprintf("{%v %voz}", c.beanType, c.ounces)
}

func (c *Coffee) BeanType() string {
	// return c.beanType
	return ""
}

func (c *Coffee) Ounces() int {
	return c.ounces
}
