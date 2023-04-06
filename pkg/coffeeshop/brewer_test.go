package coffeeshop

import (
	"coffeeshop/pkg/config"
	"coffeeshop/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBrewer_Brew(t *testing.T) {
	cfg := &config.BrewerCfg{OuncesPerSecond: 50}
	brewer := NewBrewer(cfg)
	done := make(chan *Brewer)
	tm := time.Now()
	brewer.StartBrew(model.Beans{
		BeanType:    "ethiopian",
		WeightGrams: 50,
	}, 5, func() {
		done <- brewer
	})
	<-done
	// should be +100ms
	assert.WithinRange(t, time.Now(), tm.Add(time.Millisecond*95), tm.Add(time.Millisecond*105))
}
