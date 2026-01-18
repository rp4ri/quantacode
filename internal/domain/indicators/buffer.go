package indicators

import "fmt"

// CircularBuffer maintains a fixed-size window of float64 values with O(1) updates.
type CircularBuffer struct {
	data  []float64
	size  int
	count int
	index int
	sum   float64
}

// NewCircularBuffer creates a buffer with the provided size.
func NewCircularBuffer(size int) (*CircularBuffer, error) {
	if size <= 0 {
		return nil, fmt.Errorf("buffer size must be positive")
	}
	return &CircularBuffer{data: make([]float64, size), size: size}, nil
}

// Push inserts a value, evicting the oldest when full.
func (cb *CircularBuffer) Push(value float64) {
	if cb.count == cb.size {
		old := cb.data[cb.index]
		cb.sum -= old
	} else {
		cb.count++
	}

	cb.data[cb.index] = value
	cb.sum += value
	cb.index = (cb.index + 1) % cb.size
}

// Sum returns the sum of the values currently in the buffer.
func (cb *CircularBuffer) Sum() float64 {
	return cb.sum
}

// Len returns the number of elements currently stored.
func (cb *CircularBuffer) Len() int {
	return cb.count
}

// Full returns true when the buffer has been completely filled at least once.
func (cb *CircularBuffer) Full() bool {
	return cb.count == cb.size
}

// Values returns a copy of the buffer contents in chronological order (oldest first).
func (cb *CircularBuffer) Values() []float64 {
	values := make([]float64, cb.count)
	if cb.count == 0 {
		return values
	}

	start := (cb.index - cb.count + cb.size) % cb.size
	for i := 0; i < cb.count; i++ {
		values[i] = cb.data[(start+i)%cb.size]
	}
	return values
}
