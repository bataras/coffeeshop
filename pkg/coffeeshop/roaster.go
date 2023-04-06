package coffeeshop

import (
	"coffeeshop/pkg/model"
	"coffeeshop/pkg/util"
)

type Roaster struct {
	log *util.Logger
}

func NewRoaster() *Roaster {
	return &Roaster{log: util.NewLogger("Roaster")}
}

func (r *Roaster) GetBeans(gramsNeeded int, beanType string) model.Beans {
	r.log.Infof("getbeans %v", gramsNeeded)
	return model.Beans{BeanType: beanType, WeightGrams: gramsNeeded}
}
