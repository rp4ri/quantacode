package indicators

import "fmt"

// RSI implements the Relative Strength Index using Wilder's smoothing.
type RSI struct {
	period      int
	buf         *CircularBuffer
	prev        float64
	count       int
	avgGain     float64
	avgLoss     float64
	ready       bool
	value       float64
	initialized bool
}

// NewRSI creates an RSI with the given period.
func NewRSI(period int) (*RSI, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be positive")
	}
	buf, err := NewCircularBuffer(period + 1)
	if err != nil {
		return nil, err
	}
	return &RSI{period: period, buf: buf}, nil
}

// Update ingests a price and returns the current RSI value.
// Returns 0 until sufficient data (period+1 prices) are available.
func (r *RSI) Update(price float64) float64 {
	if !r.initialized {
		r.prev = price
		r.initialized = true
		r.buf.Push(price)
		return 0
	}

	r.buf.Push(price)
	delta := price - r.prev
	var gain, loss float64
	if delta > 0 {
		gain = delta
	} else {
		loss = -delta
	}

	if r.count < r.period {
		// build initial averages via SMA over first period
		r.avgGain += gain
		r.avgLoss += loss
		r.count++
		r.prev = price
		if r.count < r.period {
			r.value = 0
			return r.value
		}
		r.avgGain /= float64(r.period)
		r.avgLoss /= float64(r.period)
		r.ready = true
	} else {
		r.avgGain = (r.avgGain*float64(r.period-1) + gain) / float64(r.period)
		r.avgLoss = (r.avgLoss*float64(r.period-1) + loss) / float64(r.period)
	}

	r.prev = price

	if !r.ready {
		r.value = 0
		return r.value
	}

	if r.avgLoss == 0 {
		r.value = 100
		return r.value
	}

	rs := r.avgGain / r.avgLoss
	r.value = 100 - (100 / (1 + rs))
	return r.value
}

// Value returns the last computed RSI.
func (r *RSI) Value() float64 {
	return r.value
}

// Period returns the configured period.
func (r *RSI) Period() int {
	return r.period
}
