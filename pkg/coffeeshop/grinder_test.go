package coffeeshop

import (
	"coffeeshop/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGrinder_Refill(t *testing.T) {
	g := NewGrinder(model.Columbian, 1000, 1000, 100, 50)

	called := 0
	err := g.Refill(func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 34,
		}
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, called)

	err = g.Refill(func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 66, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 27,
		}
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, called)

	called = 0
	err = g.Refill(func(gramsNeeded int) model.Beans {
		called++
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 27,
		}
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, called)
}

func TestGrinder_RefillWrongType(t *testing.T) {
	g := NewGrinder(model.Columbian, 1000, 1000, 100, 50)

	called := 0
	err := g.Refill(func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Ethiopian, // wrong type
			WeightGrams: 34,
		}
	})
	assert.Error(t, err)
	assert.Equal(t, 1, called)
}

func TestGrinder_RefillTime(t *testing.T) {
	g := NewGrinder(model.Columbian, 10, 50, 100, 50)

	tm := time.Now()
	called := 0
	err := g.Refill(func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 5,
		}
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	// should be +100ms
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*95), tm.Add(time.Millisecond*105))
}

func TestGrinder_Grind(t *testing.T) {
	g := NewGrinder(model.Columbian, 1000, 1000, 100, 50)

	called := 0
	beans, err := g.Grind(12, func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 34,
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 12, beans.WeightGrams)
	assert.Equal(t, model.Columbian, beans.BeanType)
	assert.Equal(t, 22, g.PercentFull())
}

func TestGrinder_GrindTime(t *testing.T) {
	g := NewGrinder(model.Columbian, 50, 50, 100, 50)

	tm := time.Now()
	called := 0
	beans, err := g.Grind(5, func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 5,
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 5, beans.WeightGrams)
	assert.Equal(t, model.Columbian, beans.BeanType)
	assert.Equal(t, 0, g.PercentFull())
	// should be +200ms (100 to refill 5g and 100 to grind 5g)
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*195), tm.Add(time.Millisecond*205))
}

func TestGrinder_GrindNotEnough(t *testing.T) {
	g := NewGrinder(model.Columbian, 1000, 1000, 100, 50)

	called := 0
	beans, err := g.Grind(12, func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Columbian,
			WeightGrams: 3, // only provide 3g
		}
	})

	assert.Error(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 0, beans.WeightGrams)
	assert.Equal(t, model.Columbian, beans.BeanType)
	assert.Equal(t, 3, g.PercentFull())
}

func TestGrinder_GrindWrongBeans(t *testing.T) {
	g := NewGrinder(model.Columbian, 1000, 1000, 100, 50)

	called := 0
	beans, err := g.Grind(12, func(gramsNeeded int) model.Beans {
		called++
		assert.Equal(t, 100, gramsNeeded)
		return model.Beans{
			BeanType:    model.Ethiopian, // wrong type
			WeightGrams: 34,
		}
	})

	assert.Error(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 0, beans.WeightGrams)
	assert.Equal(t, model.Columbian, beans.BeanType)
	assert.Equal(t, 0, g.PercentFull())
}
