package coffeeshop

import (
	"fmt"
	"math/rand"
)

type CoffeeShop struct {
	grinders                 []*Grinder
	brewers                  []*Brewer
	totalAmountUngroundBeans int
	gchan                    chan *Grinder
	bchan                    chan *Brewer
}

type Coffee struct {
	// should hold size maybe?
}

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {
	shop := CoffeeShop{
		grinders: grinders,
		brewers:  brewers,
		gchan:    make(chan *Grinder, len(grinders)),
		bchan:    make(chan *Brewer, len(brewers)),
	}

	for _, g := range grinders {
		shop.gchan <- g
	}

	for _, b := range brewers {
		shop.bchan <- b
	}

	return &shop
}

func (cs *CoffeeShop) MakeCoffee(order Order) (Coffee, error) {
	fmt.Printf("make order %v\n", order)
	// assume that we need 2 grams of beans for 1 ounce of coffee
	gramsNeededPerOunce := 2
	ungroundBeans := Beans{weightGrams: gramsNeededPerOunce * order.OuncesOfCoffeeWanted}

	grinder, ok := <-cs.gchan
	if !ok {
		return Coffee{}, fmt.Errorf("closed")
	}

	groundBeans := grinder.Grind(ungroundBeans)
	cs.gchan <- grinder

	brewer, ok := <-cs.bchan
	if !ok {
		return Coffee{}, fmt.Errorf("closed")
	}

	coffee := brewer.Brew(groundBeans)
	cs.bchan <- brewer
	return coffee, nil

	//// choose a random grinder and grind the beans
	//grinderIdx := rand.Intn(len(cs.grinders))
	//groundBeans := cs.grinders[grinderIdx].Grind(ungroundBeans)

	// NOTE: the above is for illustration purposes and does not work, because we are not considering that certain
	// grinders and brewers can be busy!

	//brewerIdx := rand.Intn(len(cs.brewers))
	//return cs.brewers[brewerIdx].Brew(groundBeans)
}

func (cs *CoffeeShop) MakeCoffeeOrg(order Order) Coffee {
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
