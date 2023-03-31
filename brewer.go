package main

import "fmt"

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	ouncesWaterPerSecond int
}

func (b *Brewer) Brew(beans Beans) Coffee {
	// assume we need 6 ounces of water for every 12 grams of beans
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)
	fmt.Printf("brew %v\n", beans)
	return Coffee{}
}
