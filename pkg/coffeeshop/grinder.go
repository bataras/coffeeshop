package coffeeshop

import (
	"fmt"
	"time"
)

type Grinder struct {
	GramsPerSecond int
}

func (g *Grinder) Grind(beans Beans) Beans {
	// how long should it take this function to complete?
	// i.e. time.Sleep(XXX)
	fmt.Printf("grind %v\n", beans)
	time.Sleep(2 * time.Second)
	return Beans{}
}
