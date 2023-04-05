package coffeeshop

import (
	"coffeeshop/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBrewer_Brew(t *testing.T) {
	brewer := NewBrewer(50)
	done := make(chan *Brewer)
	tm := time.Now()
	brewer.StartBrew(model.Beans{
		BeanType:    model.Ethiopian,
		WeightGrams: 50,
	}, 5, func() {
		done <- brewer
	})
	<-done
	// should be +100ms
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*95), tm.Add(time.Millisecond*105))
}
