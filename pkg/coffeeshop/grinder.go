package coffeeshop

import "fmt"

type Grinder struct {
	GramsPerSecond int
}

func (g *Grinder) Grind(beans Beans) Beans {
	// how long should it take this function to complete?
	// i.e. time.Sleep(XXX)
	fmt.Printf("grind %v\n", beans)
	return Beans{}
}
