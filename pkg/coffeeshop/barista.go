package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
	"fmt"
)

type Barista struct {
	shop *CoffeeShop
	log  *util.Logger
}

func NewBarista(shop *CoffeeShop) *Barista {
	return &Barista{
		shop: shop,
		log:  util.NewLogger("Barista"),
	}
}

func (b *Barista) Work() {
	beanTypes := model.BeanTypeMap()
	shop := b.shop

	go func() {
		for {
			order := shop.cashRegister.Barista()
			b.log.Infof("got order %v", order)
			if !beanTypes[order.BeanType] {
				b.log.Infof("bean type unavailable %v", order)
				order.Complete(&model.Receipt{
					Coffee: nil,
					Err:    fmt.Errorf("bean type unavailable %v", order.BeanType),
				})
			}
		}
	}()
}
