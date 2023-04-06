package coffeeshop

import (
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func cfgFor(bean string, grindGramsPerSecond int, addGramsPerSecond int, hopperSize int, refillPercentage int) *config.GrinderCfg {
	cfg := &config.GrinderCfg{
		BeanId:              "id",
		BeanCfg:             &config.BeanCfg{BeanType: bean},
		GrindGramsPerSecond: grindGramsPerSecond,
		AddGramsPerSecond:   addGramsPerSecond,
		HopperSize:          hopperSize,
		RefillPercentage:    refillPercentage,
	}
	return cfg
}

func TestGrinder_Refill(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 1000, 1000, 100, 50))

	called := 0
	err := g.TryRefill(IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		fmt.Printf("ask for %d %s\n", gramsNeeded, beanType)
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 34,
		}
	}))
	assert.NoError(t, err)
	assert.Equal(t, 1, called)

	err = g.TryRefill(IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		fmt.Printf("ask for %d %s\n", gramsNeeded, beanType)
		called++
		assert.Equal(t, 66, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 27,
		}
	}))
	assert.NoError(t, err)
	assert.Equal(t, 2, called)

	called = 0
	err = g.TryRefill(IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		fmt.Printf("ask for %d %s\n", gramsNeeded, beanType)
		called++
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 27,
		}
	}))
	assert.NoError(t, err)
	assert.Equal(t, 0, called)
}

func TestGrinder_RefillWrongType(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 1000, 1000, 100, 50))

	called := 0
	err := g.TryRefill(IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "ethiopian", // wrong type
			WeightGrams: 34,
		}
	}))
	assert.Error(t, err)
	assert.Equal(t, 1, called)
}

func TestGrinder_RefillTime(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 10, 50, 100, 50))

	tm := time.Now()
	called := 0
	err := g.TryRefill(IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 5,
		}
	}))
	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	// should be +100ms
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*95), tm.Add(time.Millisecond*105))
}

func TestGrinder_Grind(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 1000, 1000, 100, 50))

	called := 0
	beans, err := g.Grind(12, IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 34,
		}
	}))

	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 12, beans.WeightGrams)
	assert.Equal(t, "columbian", beans.BeanType)
	assert.Equal(t, 22, g.hopper.PercentFull())
	assert.True(t, g.ShouldRefill())
}

func TestGrinder_GrindTime(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 50, 50, 100, 50))

	tm := time.Now()
	called := 0
	beans, err := g.Grind(5, IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 5,
		}
	}))

	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 5, beans.WeightGrams)
	assert.Equal(t, "columbian", beans.BeanType)
	assert.Equal(t, 0, g.hopper.PercentFull())
	assert.True(t, g.ShouldRefill())
	// should be +200ms (100 to refill 5g and 100 to grind 5g)
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*195), tm.Add(time.Millisecond*205))
}

func TestGrinder_GrindNotEnough(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 1000, 1000, 100, 50))

	called := 0
	beans, err := g.Grind(12, IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "columbian",
			WeightGrams: 3, // only provide 3g
		}
	}))

	assert.Error(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 0, beans.WeightGrams)
	assert.Equal(t, "", beans.BeanType)
	assert.Equal(t, 3, g.hopper.PercentFull())
	assert.True(t, g.ShouldRefill())
}

func TestGrinder_GrindWrongBeans(t *testing.T) {
	g := NewGrinder(cfgFor("columbian", 1000, 1000, 100, 50))

	called := 0
	beans, err := g.Grind(12, IRoasterFunc(func(gramsNeeded int, beanType string) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    "ethiopian", // wrong type
			WeightGrams: 34,
		}
	}))

	assert.Error(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 0, beans.WeightGrams)
	assert.Equal(t, "", beans.BeanType)
	assert.Equal(t, 0, g.hopper.PercentFull())
	assert.True(t, g.ShouldRefill())
}
