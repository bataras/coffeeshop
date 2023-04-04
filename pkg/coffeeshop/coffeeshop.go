package coffeeshop

import (
	"coffeeshop/pkg/middleware"
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/queue"
	"coffeeshop/pkg/util"
	"fmt"
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
	cashRegister       *queue.CashRegister
	barista            *Barista
	orderQueue         chan *model.Order
	brewerDone         chan *Brewer
	grinderRefill      chan *Grinder
	beanTypes          map[model.BeanType]bool
	orderMiddleware    *middleware.Orders
	log                *util.Logger
}

const cashRegisterTimeMS int = 200
const orderQueueSize int = 4

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer) *CoffeeShop {
	cashRegister := queue.NewCashRegister(cashRegisterTimeMS)
	shop := CoffeeShop{
		extractionProfiles: NewExtractionProfiles(),
		roaster:            NewRoaster(),
		gchan:              make(chan *Grinder, len(grinders)),
		bchan:              make(chan *Brewer, len(brewers)),
		grinderRefill:      make(chan *Grinder, len(grinders)),
		brewerDone:         make(chan *Brewer, len(brewers)),
		cashRegister:       cashRegister,
		orderQueue:         make(chan *model.Order, orderQueueSize),
		orderMiddleware:    middleware.NewOrders(),
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

	shop.barista.Work()

	return &shop
}

// OrderCoffee fires off an order and returns a channel for the customer to wait on
func (cs *CoffeeShop) OrderCoffee(beanType model.BeanType, ounces int, strength model.Strength) <-chan *model.Receipt {
	rsp := make(chan *model.Receipt)
	order := model.NewOrder(rsp, cs.orderMiddleware)
	order.BeanType = beanType
	order.OuncesOfCoffeeWanted = ounces
	order.StrengthWanted = strength

	cs.cashRegister.Customer(order)

	go func(order *model.Order) {
		coffee, err := cs.makeCoffee(order)
		rsp <- &model.Receipt{
			Coffee: coffee,
			Err:    err,
		}
	}(order)

	return rsp
}

// do the work (for now)
func (cs *CoffeeShop) makeCoffee(order *model.Order) (*model.Coffee, error) {
	cs.log.Infof("make order %v\n", order)

	extractionProfile := cs.getExtractionProfile(order.StrengthWanted)
	beansNeeded := extractionProfile.GramsFromOunces(order.OuncesOfCoffeeWanted)

	// wait for a grinder
	grinder, ok := <-cs.gchan
	if !ok {
		return nil, fmt.Errorf("closed")
	}

	groundBeans, _ := grinder.Grind(beansNeeded, cs.grinderRefill, cs.roaster)
	cs.gchan <- grinder // put it back

	// wait for a brewer
	brewer, ok := <-cs.bchan
	if !ok {
		return nil, fmt.Errorf("closed")
	}

	brewer.Brew(groundBeans, order.OuncesOfCoffeeWanted, cs.brewerDone)
	brewer = <-cs.brewerDone
	coffee := brewer.GetCoffee()
	cs.bchan <- brewer // put it back
	return coffee, nil
}

func (cs *CoffeeShop) getExtractionProfile(strength model.Strength) IExtractionProfile {
	switch strength {
	default:
		fallthrough
	case model.NormalStrength:
		return cs.extractionProfiles.GetProfile(Normal)
	case model.MediumStrength:
		return cs.extractionProfiles.GetProfile(Medium)
	case model.LightStrength:
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
	return cs.brewers[brewerIdx].Brew(groundBeans, order.OuncesOfCoffeeWanted)
}
*/
