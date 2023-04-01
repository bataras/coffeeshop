package coffeeshop

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	// todo: lower case (and elsewhere like Grinder)
	OuncesWaterPerSecond int
}

// todo: possibly interact with Beans
func (b *Brewer) Brew(beans Beans, ounces int) Coffee {
	// assume we need 6 ounces of water for every 12 grams of beans
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)

	ms := time.Duration(ounces*1000/b.OuncesWaterPerSecond) * time.Millisecond
	log.Infof("brew beans %v ounces %v ms %v\n", beans, ounces, ms.Milliseconds())
	time.Sleep(ms)
	return Coffee{}
}
