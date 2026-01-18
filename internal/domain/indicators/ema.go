package indicators

import "fmt"

// EMA implements an exponential moving average over a fixed period.
type EMA struct {
	period      int
	buf         *CircularBuffer
	multiplier  float64
	value       float64
	initialized bool
}

// NewEMA creates an EMA with the given period.
func NewEMA(period int) (*EMA, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be positive")
	}
	buf, err := NewCircularBuffer(period)
	if err != nil {
		return nil, err
	}
	return &EMA{
		period:     period,
		buf:        buf,
		multiplier: 2.0 / float64(period+1),
	}, nil
}

// Update ingests a price and returns the current EMA value.
// Returns 0 until the buffer is filled; initializes using SMA of the first period prices.
func (e *EMA) Update(price float64) float64 {
	e.buf.Push(price)

	if !e.initialized {
		if e.buf.Len() < e.period {
			e.value = 0
			return e.value
		}
		// initialize using SMA of first period
		e.value = e.buf.Sum() / float64(e.period)
		e.initialized = true
		return e.value
	}

	e.value = (price-e.value)*e.multiplier + e.value
	return e.value
}

// Value returns the last computed EMA.
func (e *EMA) Value() float64 {
	return e.value
}

// Period returns the configured period.
func (e *EMA) Period() int {
	return e.period
}
