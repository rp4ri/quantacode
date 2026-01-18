package indicators

import "fmt"

// SMA implements a simple moving average over a fixed period.
type SMA struct {
	period int
	buf    *CircularBuffer
	value  float64
}

// NewSMA creates an SMA with the given period.
func NewSMA(period int) (*SMA, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be positive")
	}
	buf, err := NewCircularBuffer(period)
	if err != nil {
		return nil, err
	}
	return &SMA{period: period, buf: buf}, nil
}

// Update ingests a price and returns the current SMA value. Returns 0 until enough data is collected.
func (s *SMA) Update(price float64) float64 {
	s.buf.Push(price)
	if s.buf.Len() < s.period {
		s.value = 0
		return s.value
	}
	s.value = s.buf.Sum() / float64(s.period)
	return s.value
}

// Value returns the last computed SMA.
func (s *SMA) Value() float64 {
	return s.value
}

// Period returns the configured period.
func (s *SMA) Period() int {
	return s.period
}
