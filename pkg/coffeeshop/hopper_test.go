package coffeeshop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHopper(t *testing.T) {
	h := NewHopper(-1)
	assert.Equal(t, 0, h.Count())
	assert.Equal(t, 0, h.Size())
	assert.Equal(t, 0, h.SpaceAvailable())
	assert.Equal(t, 0, h.AddBeans(5))
	assert.Equal(t, 0, h.Count())
	assert.Equal(t, 0, h.TakeBeans(5))
}

func TestHopper2(t *testing.T) {
	h := NewHopper(200)
	assert.Equal(t, 0, h.Count())
	assert.Equal(t, 200, h.Size())
	assert.Equal(t, 200, h.SpaceAvailable())
	assert.Equal(t, 5, h.AddBeans(5))
	assert.Equal(t, 5, h.TakeBeans(50))
	assert.Equal(t, 200, h.AddBeans(500))
	assert.Equal(t, 0, h.TakeBeans(-1))
	assert.Equal(t, 0, h.AddBeans(-1))
	assert.Equal(t, 0, h.SpaceAvailable())
	assert.Equal(t, 30, h.TakeBeans(30))
	assert.Equal(t, 30, h.SpaceAvailable())
	assert.Equal(t, 170, h.Count())
	assert.Equal(t, 30, h.TakeBeans(30))
	assert.Equal(t, 60, h.SpaceAvailable())
	assert.Equal(t, 70, h.PercentFull())
	assert.Equal(t, 15, h.TakeBeans(15))
	assert.Equal(t, 75, h.SpaceAvailable())
	assert.Equal(t, 62, h.PercentFull())
}
