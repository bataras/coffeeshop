package util

import "time"

// RateToDuration convert amount and rate into a duration
func RateToDuration(amount, perSecond int) time.Duration {
	if perSecond == 0 {
		return time.Duration(0)
	}
	return time.Duration(amount*1000/perSecond) * time.Millisecond
}

// Rate simple struct for getting a duration from a rate and amount
type Rate struct {
	perSecond int
}

func (r *Rate) SetPerSecond(rate int) {
	r.perSecond = rate
}

func (r *Rate) Duration(amount int) time.Duration {
	return RateToDuration(amount, r.perSecond)
}
