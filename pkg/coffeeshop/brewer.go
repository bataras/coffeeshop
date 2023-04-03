package coffeeshop

import (
	"coffeeshop/pkg/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	ouncesPerSecond util.Rate
}

func NewBrewer(ouncesPerSecond int) *Brewer {
	val := &Brewer{}
	val.ouncesPerSecond.SetPerSecond(ouncesPerSecond)
	return val
}

// Brew todo: possibly interact with Beans
func (b *Brewer) Brew(beans Beans, ounces int) *Coffee {
	// assume we need 6 ounces of water for every 12 grams of beans
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)

	ms := b.ouncesPerSecond.Duration(ounces)
	log.Infof("brew beans %v ounces %v ms %v\n", beans, ounces, ms.Milliseconds())
	time.Sleep(ms)
	return NewCoffee(beans.beanType, ounces)
}
