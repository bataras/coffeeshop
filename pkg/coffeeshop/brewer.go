package coffeeshop

import (
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"time"
)

type BrewerDoneCB func()

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	ouncesPerSecond util.Rate
	coffee          *model.Coffee
	log             *util.Logger
}

func NewBrewer(cfg *config.BrewerCfg) *Brewer {
	val := &Brewer{
		log: util.NewLogger("Brewer"),
	}
	val.ouncesPerSecond.SetPerSecond(cfg.OuncesPerSecond)
	return val
}

func (b *Brewer) StartBrew(beans model.Beans, ounces int, whenDone BrewerDoneCB) {
	// assume we need 6 ounces of water for every 12 grams of beans
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)

	go func() {
		ms := b.ouncesPerSecond.Duration(ounces)
		b.log.Infof("brew beans %v ounces %v ms %v", beans, ounces, ms.Milliseconds())
		time.Sleep(ms)
		b.coffee = model.NewCoffee(beans.BeanType, ounces)
		if whenDone != nil {
			whenDone() // let the caller do something later
		}
	}()
}

func (b *Brewer) GetCoffee() *model.Coffee {
	return b.coffee
}
