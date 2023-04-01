package main

import (
	"coffeeshop/pkg/coffeeshop"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func init() {
	// todo: config logging
	log.SetOutput(os.Stdout)
	//log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMicro,
	})
}

// todo: add app config
func main() {
	// Premise: we want to model a coffee shop. An order comes in, and then with a limited amount of grinders and
	// brewers (each of which can be "busy"): we must grind unground beans, take the resulting ground beans, and then
	// brew them into liquid coffee. We need to coordinate the work when grinders and/or brewers are busy doing work
	// already. What Go datastructure(s) might help us coordinate the steps: order -> grinder -> brewer -> coffee?
	//
	// Some struct types and their functions need to be filled in properly. It may be helpful to finish the
	// Grinder impl, and then Brewer impl each, and then see how things all fit together inside CoffeeShop afterwards.

	g1 := &coffeeshop.Grinder{GramsPerSecond: 5}
	g2 := &coffeeshop.Grinder{GramsPerSecond: 3}
	g3 := &coffeeshop.Grinder{GramsPerSecond: 12}

	b1 := &coffeeshop.Brewer{OuncesWaterPerSecond: 2}
	b2 := &coffeeshop.Brewer{OuncesWaterPerSecond: 5}

	cs := coffeeshop.NewCoffeeShop([]*coffeeshop.Grinder{g1, g2, g3}, []*coffeeshop.Brewer{b1, b2})

	var wg sync.WaitGroup
	numCustomers := 10
	for i := 0; i < numCustomers; i++ {
		// in parallel, all at once, make calls to MakeCoffee
		wg.Add(1)
		go func() {
			coffee, err := cs.MakeCoffee(coffeeshop.Order{OuncesOfCoffeeWanted: 12})
			log.Infof("made %v err %v\n", coffee, err)
			wg.Done()
		}()
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
