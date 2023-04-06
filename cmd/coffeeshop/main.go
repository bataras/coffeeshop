package main

import (
	"coffeeshop/pkg/coffeeshop"
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/util"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"sync"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	// log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMicro,
	})
}

// todo: error handling
// todo: context shutown system-wide
func main() {
	log := util.NewLogger("Main")

	cfg, err := config.Load("coffeeshop.yaml")
	if err != nil {
		log.Errorf("config problem: %v", err)
		os.Exit(1)
	}

	// Premise: we want to model a coffee shop. An order comes in, and then with a limited amount of grinders and
	// brewers (each of which can be "busy"): we must grind unground beans, take the resulting ground beans, and then
	// brew them into liquid coffee. We need to coordinate the work when grinders and/or brewers are busy doing work
	// already. What Go datastructure(s) might help us coordinate the steps: order -> grinder -> brewer -> coffee?
	//
	// Some struct types and their functions need to be filled in properly. It may be helpful to finish the
	// Grinder impl, and then Brewer impl each, and then see how things all fit together inside CoffeeShop afterwards.

	cs := coffeeshop.NewCoffeeShop(cfg)

	var wg sync.WaitGroup
	doOrder := func(bean string, ounces int, strength coffeeshop.Strength) {
		receipt := cs.OrderCoffee(bean, ounces, strength)
		wg.Add(1)
		go func() {
			coffee, ok := <-receipt
			if !ok {
				log.Infof("order closed")
			} else {
				if coffee.Err != nil {
					log.Errorf("order handling error %v", coffee.Err)
				} else {
					log.Infof("made Order#: %d Coffee: %s", coffee.OrderNumber, coffee.Coffee.String())
				}
			}
			wg.Done()
		}()
	}

	var beanTypes []string
	for bt, _ := range cfg.BeanTypes() {
		beanTypes = append(beanTypes, bt)
	}

	strengths := coffeeshop.OrderStrengths()
	numCustomers := cfg.Shop.CustomerCount
	for i := 0; i < numCustomers; i++ {
		randomIndex := rand.Intn(len(beanTypes))
		randomSize := rand.Intn(8) + 8
		randomStrength := strengths[rand.Intn(len(strengths))]
		doOrder(beanTypes[randomIndex], randomSize, randomStrength)
	}
	wg.Wait()

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
