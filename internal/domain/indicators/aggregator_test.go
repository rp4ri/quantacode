package indicators

import (
	"testing"
)

func TestNewAggregator(t *testing.T) {
	agg, err := NewAggregator(14, 14, 14)
	if err != nil {
		t.Fatalf("NewAggregator() error = %v", err)
	}
	if agg == nil {
		t.Fatal("NewAggregator() returned nil")
	}
}

func TestNewAggregatorInvalidPeriod(t *testing.T) {
	_, err := NewAggregator(0, 14, 14)
	if err == nil {
		t.Error("NewAggregator() should fail with period 0")
	}

	_, err = NewAggregator(14, -1, 14)
	if err == nil {
		t.Error("NewAggregator() should fail with negative period")
	}
}

func TestAggregatorUpdate(t *testing.T) {
	agg, _ := NewAggregator(14, 14, 14)

	prices := []float64{100, 102, 101, 103, 105, 104, 106, 108, 107, 109, 110, 112, 111, 113, 115}

	var lastVals AggregatedValues
	for _, p := range prices {
		lastVals = agg.Update(p)
	}

	if lastVals.SMA == 0 {
		t.Error("SMA should not be 0 after 15 updates")
	}
	if lastVals.EMA == 0 {
		t.Error("EMA should not be 0 after 15 updates")
	}
}

func TestAggregatorHistory(t *testing.T) {
	agg, _ := NewAggregator(14, 14, 14)

	for i := 0; i < 35; i++ {
		agg.Update(float64(100 + i))
	}

	history := agg.History()

	if len(history.RSI) != HistorySize {
		t.Errorf("History RSI length = %d, want %d", len(history.RSI), HistorySize)
	}
	if len(history.SMA) != HistorySize {
		t.Errorf("History SMA length = %d, want %d", len(history.SMA), HistorySize)
	}
	if len(history.EMA) != HistorySize {
		t.Errorf("History EMA length = %d, want %d", len(history.EMA), HistorySize)
	}
	if len(history.Prices) != HistorySize {
		t.Errorf("History Prices length = %d, want %d", len(history.Prices), HistorySize)
	}
}

func TestAggregatorHistoryIsCopy(t *testing.T) {
	agg, _ := NewAggregator(14, 14, 14)

	for i := 0; i < 20; i++ {
		agg.Update(float64(100 + i))
	}

	history1 := agg.History()
	originalLen := len(history1.RSI)

	agg.Update(200)

	history2 := agg.History()

	if len(history1.RSI) != originalLen {
		t.Error("History() should return a copy, not a reference")
	}
	if len(history2.RSI) != originalLen+1 {
		t.Error("New history should have one more element")
	}
}

func TestAppendWithLimit(t *testing.T) {
	tests := []struct {
		name     string
		slice    []float64
		val      float64
		limit    int
		wantLen  int
		wantLast float64
	}{
		{
			name:     "append to empty",
			slice:    nil,
			val:      1.0,
			limit:    5,
			wantLen:  1,
			wantLast: 1.0,
		},
		{
			name:     "append under limit",
			slice:    []float64{1, 2, 3},
			val:      4.0,
			limit:    5,
			wantLen:  4,
			wantLast: 4.0,
		},
		{
			name:     "append at limit",
			slice:    []float64{1, 2, 3, 4, 5},
			val:      6.0,
			limit:    5,
			wantLen:  5,
			wantLast: 6.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendWithLimit(tt.slice, tt.val, tt.limit)
			if len(got) != tt.wantLen {
				t.Errorf("appendWithLimit() len = %d, want %d", len(got), tt.wantLen)
			}
			if got[len(got)-1] != tt.wantLast {
				t.Errorf("appendWithLimit() last = %v, want %v", got[len(got)-1], tt.wantLast)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"single value", []int{5}, 5},
		{"multiple values", []int{1, 5, 3}, 5},
		{"negative values", []int{-1, -5, -3}, -1},
		{"empty", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxInt(tt.values...); got != tt.want {
				t.Errorf("maxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
