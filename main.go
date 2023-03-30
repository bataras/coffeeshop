package main

import (
	"math/rand"
)

type Beans struct {
	weightGrams int
	// indicate some state change? create a new type?
}

type Grinder struct {
	gramsPerSecond int
}

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	ouncesWaterPerSecond int
}

type Order struct {
	ouncesOfCoffeeWanted int
}

type CoffeeShop struct {
	grinders                 []*Grinder
	brewers                  []*Brewer
	totalAmountUngroundBeans int
}

type Coffee struct {
	// should hold size maybe?
}

func (g *Grinder) Grind(beans Beans) Beans {
	// how long should it take this function to complete?
	// i.e. time.Sleep(XXX)
	return Beans{}
}

func (b *Brewer) Brew(beans Beans) Coffee {
	// assume we need 6 ounces of water for every 12 grams of beans
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)
	return Coffee{}
}

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {
	return &CoffeeShop{grinders: grinders, brewers: brewers}
}

func (cs *CoffeeShop) MakeCoffee(order Order) Coffee {
	// assume that we need 2 grams of beans for 1 ounce of coffee
	gramsNeededPerOunce := 2
	ungroundBeans := Beans{weightGrams: gramsNeededPerOunce * order.ouncesOfCoffeeWanted}
	// choose a random grinder and grind the beans
	grinderIdx := rand.Intn(len(cs.grinders))
	groundBeans := cs.grinders[grinderIdx].Grind(ungroundBeans)

	// NOTE: the above is for illustration purposes and does not work, because we are not considering that certain
	// grinders and brewers can be busy!

	brewerIdx := rand.Intn(len(cs.brewers))
	return cs.brewers[brewerIdx].Brew(groundBeans)
}

func main() {
	// Premise: we want to model a coffee shop. An order comes in, and then with a limited amount of grinders and
	// brewers (each of which can be "busy"): we must grind unground beans, take the resulting ground beans, and then
	// brew them into liquid coffee. We need to coordinate the work when grinders and/or brewers are busy doing work
	// already. What Go datastructure(s) might help us coordinate the steps: order -> grinder -> brewer -> coffee?
	//
	// Some of the struct types and their functions need to be filled in properly. It may be helpful to finish the
	// Grinder impl, and then Brewer impl each, and then see how things all fit together inside CoffeeShop afterwards.

	//b := Beans{weightGrams: 10}
	g1 := &Grinder{gramsPerSecond: 5}
	g2 := &Grinder{gramsPerSecond: 3}
	g3 := &Grinder{gramsPerSecond: 12}

	b1 := &Brewer{ouncesWaterPerSecond: 100}
	b2 := &Brewer{ouncesWaterPerSecond: 25}

	cs := NewCoffeeShop([]*Grinder{g1, g2, g3}, []*Brewer{b1, b2})

	numCustomers := 10
	for i := 0; i < numCustomers; i++ {
		// in parallel, all at once, make calls to MakeCoffee
		go func() {
			cs.MakeCoffee(Order{ouncesOfCoffeeWanted: 12})
		}()
	}
	// Issues with the above
	// 1. Assumes that we have unlimited amounts of grinders and brewers.
	//		- How do we build in logic that takes into account that a given Grinder or Brewer is busy?
	// 2. Does not take into account that brewers must be used after grinders are done.
	// 		- Making a coffee needs to be done sequentially: find an open grinder, grind the beans, find an open brewer,
	//		  brew the ground beans into coffee.
	// 3. A lot of assumptions (i.e. 2 grams needed for 1 ounce of coffee) are left as comments in the code.
	// 		- How can we make these assumptions configurable, so that our coffee shop can serve let's say different
	//		  strengths of coffee via the Order that is placed (i.e. 5 grams of beans to make 1 ounce of coffee)?
}
