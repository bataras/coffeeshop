package coffeeshop

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type CoffeeShop struct {
	// todo: add baristas and possibly "floor space" to regulate the max # of baristas
	// that can brew. allow baristas to occasionally go on break and clean tables and empty a full waste hopper

	// todo: add hoppers and a waste bucket/hopper
	// todo: add multiple bean types (and hoppers)

	// baristas grab jobs: orders, filling hoppers, cleaning tables, taking breaks
	// baristas grab resources for exclusive use: hoppers, grinders, brewers
	extractionProfiles       IExtractionProfiles
	beanHopper               Hopper
	grinders                 []*Grinder
	brewers                  []*Brewer
	totalAmountUngroundBeans int
	gchan                    chan *Grinder
	bchan                    chan *Brewer
}

type Coffee struct {
	// todo: should hold size maybe?
}

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {
	shop := CoffeeShop{
		extractionProfiles: NewExtractionProfiles(),
		grinders:           grinders,
		brewers:            brewers,
		gchan:              make(chan *Grinder, len(grinders)),
		bchan:              make(chan *Brewer, len(brewers)),
	}

	for _, g := range grinders {
		shop.gchan <- g
	}

	for _, b := range brewers {
		shop.bchan <- b
	}

	shop.beanHopper.AddBeans(100)

	return &shop
}

func (cs *CoffeeShop) MakeCoffee(order Order) (Coffee, error) {
	log.Infof("make order %v\n", order)

	extractionProfile := cs.getExtractionProfile(order.StrengthWanted)
	beansNeeded := extractionProfile.GramsFromOunces(order.OuncesOfCoffeeWanted)
	err := cs.beanHopper.TakeBeans(beansNeeded)
	if err != nil {
		return Coffee{}, err
	}
	ungroundBeans := Beans{weightGrams: extractionProfile.GramsFromOunces(order.OuncesOfCoffeeWanted)}

	// wait for a grinder
	grinder, ok := <-cs.gchan
	if !ok {
		return Coffee{}, fmt.Errorf("closed")
	}

	groundBeans := grinder.Grind(ungroundBeans)
	cs.gchan <- grinder // put it back

	// wait for a brewer
	brewer, ok := <-cs.bchan
	if !ok {
		return Coffee{}, fmt.Errorf("closed")
	}

	coffee := brewer.Brew(groundBeans, order.OuncesOfCoffeeWanted)
	cs.bchan <- brewer // put it back
	return coffee, nil
}

func (cs *CoffeeShop) getExtractionProfile(strength Strength) IExtractionProfile {
	switch strength {
	default:
		fallthrough
	case NormalStrength:
		return cs.extractionProfiles.GetProfile(Normal)
	case MediumStrength:
		return cs.extractionProfiles.GetProfile(Medium)
	case LightStrength:
		return cs.extractionProfiles.GetProfile(Light)
	}
}

/*
func (cs *CoffeeShop) MakeCoffeeOrg(order Order) Coffee {
	log.Infof("make order %v\n", order)
	// assume that we need 2 grams of beans for 1 ounce of coffee
	// todo: make configurable
	gramsNeededPerOunce := 2
	ungroundBeans := Beans{weightGrams: gramsNeededPerOunce * order.OuncesOfCoffeeWanted}
	// choose a random grinder and grind the beans
	grinderIdx := rand.Intn(len(cs.grinders))
	groundBeans := cs.grinders[grinderIdx].Grind(ungroundBeans)

	// NOTE: the above is for illustration purposes and does not work, because we are not considering that certain
	// grinders and brewers can be busy!

	brewerIdx := rand.Intn(len(cs.brewers))
	return cs.brewers[brewerIdx].Brew(groundBeans, order.OuncesOfCoffeeWanted)
}
*/
