package coffeeshop

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBrewer_Brew(t *testing.T) {
	brewer := NewBrewer(50)
	tm := time.Now()
	brewer.Brew(Beans{
		beanType:    Ethiopian,
		weightGrams: 50,
	}, 5)
	// should be +100ms
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*95), tm.Add(time.Millisecond*105))
}
