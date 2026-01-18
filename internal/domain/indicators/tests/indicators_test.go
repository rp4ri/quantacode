package indicators_test

import (
	"math"
	"testing"

	"github.com/yourusername/quantacode/internal/domain/indicators"
)

func TestSMA(t *testing.T) {
	tests := []struct {
		name     string
		prices   []float64
		period   int
		expected float64
	}{
		{"insufficient data returns zero", []float64{1, 2}, 3, 0},
		{"exact period average", []float64{1, 2, 3}, 3, 2},
		{"wrap updates oldest", []float64{1, 2, 3, 4}, 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sma, err := indicators.NewSMA(tt.period)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var got float64
			for _, price := range tt.prices {
				got = sma.Update(price)
			}
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEMA(t *testing.T) {
	ema, err := indicators.NewEMA(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prices := []float64{10, 20, 30, 40}
	var got float64
	for _, price := range prices {
		got = ema.Update(price)
	}

	// period=3 => multiplier=0.5; initial EMA seeded by SMA=20; next value should be 30.
	if math.Abs(got-30) > 1e-9 {
		t.Fatalf("EMA got %v, want 30", got)
	}
}

func TestRSI(t *testing.T) {
	tests := []struct {
		name     string
		prices   []float64
		period   int
		expected float64
	}{
		{"all gains yields 100", []float64{1, 2, 3, 4}, 3, 100},
		{"mixed gains and losses", []float64{1, 2, 1, 2, 1}, 3, 44.44},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsi, err := indicators.NewRSI(tt.period)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var got float64
			for _, price := range tt.prices {
				got = rsi.Update(price)
			}
			if math.Abs(got-tt.expected) > 0.2 { // allow small float tolerance
				t.Fatalf("got %v, want ~%v", got, tt.expected)
			}
		})
	}
}

func TestAggregator(t *testing.T) {
	agg, err := indicators.NewAggregator(3, 3, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prices := []float64{1, 2, 3, 4}
	var last indicators.AggregatedValues
	for _, price := range prices {
		last = agg.Update(price)
	}

	current := agg.Values()
	if last != current {
		t.Fatalf("aggregated values mismatch: got %v, want %v", current, last)
	}

	if math.Abs(last.SMA-3) > 1e-9 {
		t.Fatalf("SMA got %v, want 3", last.SMA)
	}
}
