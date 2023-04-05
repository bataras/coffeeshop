package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
)

type CoffeeShop struct {
	// todo: add baristas and possibly "floor space" to regulate the max # of baristas
	// that can brew. allow baristas to occasionally go on break and clean tables and empty a full waste hopper

	// todo: add hoppers and a waste bucket/hopper
	// todo: add multiple bean types (and hoppers)

	// baristas grab jobs: orders, filling hoppers, cleaning tables, taking breaks
	// baristas grab resources for exclusive use: hoppers, grinders, brewers
	extractionProfiles IExtractionProfiles
	roaster            *Roaster
	gchan              chan *Grinder
	bchan              chan *Brewer
	cashRegister       *CashRegister
	barista            *Barista
	orderQueue         *util.PriorityWaitQueue[*Order]
	brewerDone         chan *Order
	grinderRefill      chan *Grinder
	beanTypes          map[model.BeanType]bool
	orderObserver      IOrderObserver
	log                *util.Logger
}

const cashRegisterTimeMS int = 200

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {

	cashRegister := NewCashRegister(cashRegisterTimeMS)
	shop := CoffeeShop{
		extractionProfiles: NewExtractionProfiles(),
		roaster:            NewRoaster(),
		gchan:              make(chan *Grinder, len(grinders)),
		bchan:              make(chan *Brewer, len(brewers)),
		grinderRefill:      make(chan *Grinder, len(grinders)),
		brewerDone:         make(chan *Order, len(brewers)),
		cashRegister:       cashRegister,
		orderQueue:         util.NewPriorityWaitQueue[*Order](),
		orderObserver:      NewOrderObserver(),
		log:                util.NewLogger("Shop"),
	}
	shop.barista = NewBarista(&shop) // todo: allow multiple baristas

	// todo build brewers/grinders from config and assign done channels here
	for _, g := range grinders {
		shop.gchan <- g
	}

	for _, b := range brewers {
		shop.bchan <- b
	}

	shop.barista.StartWork()

	return &shop
}

// OrderCoffee fires off an order and returns a channel for the customer to wait on
func (cs *CoffeeShop) OrderCoffee(beanType model.BeanType, ounces int, strength Strength) <-chan *model.Receipt {
	rsp := make(chan *model.Receipt)
	order := NewOrder(rsp, cs.orderObserver)
	order.BeanType = beanType
	order.OuncesOfCoffeeWanted = ounces
	order.StrengthWanted = strength

	cs.cashRegister.Customer(order)

	return rsp
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
	cs.log.Infof("make order %v\n", order)
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
	return cs.brewers[brewerIdx].StartBrew(groundBeans, order.OuncesOfCoffeeWanted)
}
*/
