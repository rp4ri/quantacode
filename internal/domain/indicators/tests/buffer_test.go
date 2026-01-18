package indicators_test

import (
	"testing"

	"github.com/yourusername/quantacode/internal/domain/indicators"
)

func TestCircularBufferPushAndWrap(t *testing.T) {
	buf, err := indicators.NewCircularBuffer(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf.Push(1)
	buf.Push(2)
	buf.Push(3)

	if buf.Sum() != 6 {
		t.Fatalf("sum after initial fill = %v, want 6", buf.Sum())
	}
	if buf.Len() != 3 {
		t.Fatalf("len after initial fill = %v, want 3", buf.Len())
	}

	buf.Push(4) // wrap

	if buf.Sum() != 9 {
		t.Fatalf("sum after wrap = %v, want 9", buf.Sum())
	}
	wantValues := []float64{2, 3, 4}
	values := buf.Values()
	if len(values) != len(wantValues) {
		t.Fatalf("values length = %d, want %d", len(values), len(wantValues))
	}
	for i, v := range wantValues {
		if values[i] != v {
			t.Fatalf("values[%d] = %v, want %v", i, values[i], v)
		}
	}
}
