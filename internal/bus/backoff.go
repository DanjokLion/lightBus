package bus

import (
	"math"
	"time"
)

type Backoff interface {
	NextDelay(attempt int) time.Duration
}

type ExpBackoff struct {
	Base	time.Duration
	Factor	float64
	Max 	time.Duration
}

func (b ExpBackoff) NextDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return b.Base
	}
	dur := float64(b.Base) * math.Pow(b.Factor, float64(attempt - 1))
	if max := float64(b.Max); dur > max {
		dur = max
	}
	return time.Duration(dur)

}