package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
)

type CoffeeShop struct {
	extractionProfiles IExtractionProfiles
	roaster            *Roaster
	grinders           chan *Grinder
	brewers            chan *Brewer
	cashRegister       *CashRegister
	orderPipeDepth     chan bool
	orderQueue         *util.PriorityWaitQueue[*Order]
	brewerDone         *util.PriorityWaitQueue[*Order]
	grinderRefill      chan *Grinder
	beanTypes          map[model.BeanType]bool
	log                *util.Logger
}

const cashRegisterTimeMS int = 200

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer, baristas int) *CoffeeShop {

	cashRegister := NewCashRegister(cashRegisterTimeMS)
	orderPipeDepth := len(grinders) + len(brewers) + 2 // max orders being handled in the shop
	shop := CoffeeShop{
		extractionProfiles: NewExtractionProfiles(),
		roaster:            NewRoaster(),
		grinders:           make(chan *Grinder, len(grinders)), // todo: map of grinders
		brewers:            make(chan *Brewer, len(brewers)),
		grinderRefill:      make(chan *Grinder, len(grinders)),
		brewerDone:         util.NewPriorityWaitQueue[*Order](),
		cashRegister:       cashRegister,                    // todo: more than one cash register
		orderPipeDepth:     make(chan bool, orderPipeDepth), // back pressure orders
		orderQueue:         util.NewPriorityWaitQueue[*Order](),
		log:                util.NewLogger("Shop"),
	}

	// todo build brewers/grinders from config and assign done channels here
	for _, g := range grinders {
		shop.grinders <- g
	}

	for _, b := range brewers {
		shop.brewers <- b
	}

	// fire off the baristas
	for i := 0; i < baristas; i++ {
		barista := NewBarista(&shop)
		barista.StartWork()
	}

	return &shop
}

// OrderCoffee fires off an order and returns a channel for the customer to wait on
func (cs *CoffeeShop) OrderCoffee(beanType model.BeanType, ounces int, strength Strength) <-chan *model.Receipt {
	cs.orderPipeDepth <- true

	rsp := make(chan *model.Receipt)
	order := NewOrder(rsp, func() {
		<-cs.orderPipeDepth
	})
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
	cs.log.Infof("make order %v", order)
	// assume that we need 2 grams of beans for 1 ounce of coffee
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
