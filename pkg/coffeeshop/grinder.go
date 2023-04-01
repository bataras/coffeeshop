package coffeeshop

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Grinder struct {
	GramsPerSecond int
}

func (g *Grinder) Grind(beans Beans) Beans {
	// how long should it take this function to complete?
	// i.e. time.Sleep(XXX)
	ms := time.Duration(beans.weightGrams*1000/g.GramsPerSecond) * time.Millisecond
	log.Infof("grind beans %v ms %v\n", beans, ms.Milliseconds())
	time.Sleep(ms)
	return Beans{weightGrams: beans.weightGrams}
}
