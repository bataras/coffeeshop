package coffeeshop

import (
	"fmt"
	"math/rand"
)

type CoffeeShop struct {
	grinders                 []*Grinder
	brewers                  []*Brewer
	totalAmountUngroundBeans int
}

type Coffee struct {
	// should hold size maybe?
}

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {
	return &CoffeeShop{grinders: grinders, brewers: brewers}
}

func (cs *CoffeeShop) MakeCoffee(order Order) Coffee {
	fmt.Printf("make order %v\n", order)
	// assume that we need 2 grams of beans for 1 ounce of coffee
	gramsNeededPerOunce := 2
	ungroundBeans := Beans{weightGrams: gramsNeededPerOunce * order.OuncesOfCoffeeWanted}
	// choose a random grinder and grind the beans
	grinderIdx := rand.Intn(len(cs.grinders))
	groundBeans := cs.grinders[grinderIdx].Grind(ungroundBeans)

	// NOTE: the above is for illustration purposes and does not work, because we are not considering that certain
	// grinders and brewers can be busy!

	brewerIdx := rand.Intn(len(cs.brewers))
	return cs.brewers[brewerIdx].Brew(groundBeans)
}
