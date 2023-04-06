package coffeeshop

import (
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"strings"
)

type CoffeeShop struct {
	extractionProfiles IExtractionProfiles
	roaster            *Roaster
	grinders           *GrinderPool
	brewers            chan *Brewer
	cashRegister       *CashRegister
	orderPipeDepth     chan bool
	orderQueue         *util.PriorityWaitQueue[*Order]
	brewerDone         *util.PriorityWaitQueue[*Order]
	grinderRefill      chan *Grinder
	beanTypes          map[string]bool
	log                *util.Logger
}

func NewCoffeeShop(cfg *config.Config) *CoffeeShop {

	cashRegister := NewCashRegister(cfg.Shop.CashRegisterTimeMS)
	orderPipeDepth := len(cfg.Grinders) + len(cfg.Brewers) + 2 // max orders being handled in the shop
	shop := CoffeeShop{
		extractionProfiles: NewExtractionProfiles(),
		roaster:            NewRoaster(),
		grinders:           NewGrinderPool(len(cfg.Grinders)),
		brewers:            make(chan *Brewer, len(cfg.Brewers)),
		grinderRefill:      make(chan *Grinder, len(cfg.Grinders)),
		orderQueue:         util.NewPriorityWaitQueue[*Order](),
		brewerDone:         util.NewPriorityWaitQueue[*Order](),
		cashRegister:       cashRegister,                    // todo: more than one cash register
		orderPipeDepth:     make(chan bool, orderPipeDepth), // back pressure orders
		beanTypes:          cfg.BeanTypes(),
		log:                util.NewLogger("Shop"),
	}

	for _, gcfg := range cfg.Grinders {
		grinder := NewGrinder(gcfg)
		shop.grinders.Put(grinder)
	}

	for _, bcfg := range cfg.Brewers {
		brewer := NewBrewer(bcfg)
		shop.brewers <- brewer
	}

	// fire off the baristas
	for i := 0; i < cfg.Shop.BaristaCount; i++ {
		barista := NewBarista(&shop)
		barista.StartWork()
	}

	return &shop
}

// OrderCoffee fires off an order and returns a channel for the customer to wait on
func (cs *CoffeeShop) OrderCoffee(beanType string, ounces int, strength Strength) <-chan *model.Receipt {
	cs.orderPipeDepth <- true

	rsp := make(chan *model.Receipt, 1) // cap 1: pushing the receipt should not block the barista
	order := NewOrder(rsp, func() {
		<-cs.orderPipeDepth
	})
	order.BeanType = strings.ToLower(beanType)
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
