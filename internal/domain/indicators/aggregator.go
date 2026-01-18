package indicators

const (
	// HistorySize defines how many candles of indicator history to keep
	HistorySize = 30
)

// AggregatedValues contains the latest values for all indicators.
type AggregatedValues struct {
	RSI float64
	SMA float64
	EMA float64
}

// IndicatorHistory contains historical values for indicators
type IndicatorHistory struct {
	RSI    []float64
	SMA    []float64
	EMA    []float64
	Prices []float64
}

// Aggregator coordinates price updates across indicators and tracks recent prices.
type Aggregator struct {
	prices     *CircularBuffer
	rsi        *RSI
	sma        *SMA
	ema        *EMA
	last       AggregatedValues
	rsiHistory []float64
	smaHistory []float64
	emaHistory []float64
	priceHistory []float64
}

// NewAggregator constructs an Aggregator with the provided indicator periods.
func NewAggregator(rsiPeriod, smaPeriod, emaPeriod int) (*Aggregator, error) {
	maxPeriod := maxInt(rsiPeriod+1, smaPeriod, emaPeriod)
	prices, err := NewCircularBuffer(maxPeriod)
	if err != nil {
		return nil, err
	}

	rsi, err := NewRSI(rsiPeriod)
	if err != nil {
		return nil, err
	}

	sma, err := NewSMA(smaPeriod)
	if err != nil {
		return nil, err
	}

	ema, err := NewEMA(emaPeriod)
	if err != nil {
		return nil, err
	}

	return &Aggregator{
		prices: prices,
		rsi:    rsi,
		sma:    sma,
		ema:    ema,
	}, nil
}

// Update ingests a price, updates all indicators, and returns aggregated values.
func (a *Aggregator) Update(price float64) AggregatedValues {
	a.prices.Push(price)
	rsiVal := a.rsi.Update(price)
	smaVal := a.sma.Update(price)
	emaVal := a.ema.Update(price)
	
	a.last = AggregatedValues{
		RSI: rsiVal,
		SMA: smaVal,
		EMA: emaVal,
	}
	
	a.priceHistory = appendWithLimit(a.priceHistory, price, HistorySize)
	a.rsiHistory = appendWithLimit(a.rsiHistory, rsiVal, HistorySize)
	a.smaHistory = appendWithLimit(a.smaHistory, smaVal, HistorySize)
	a.emaHistory = appendWithLimit(a.emaHistory, emaVal, HistorySize)
	
	return a.last
}

// Values returns the most recently computed aggregate values.
func (a *Aggregator) Values() AggregatedValues {
	return a.last
}

// History returns the historical values for all indicators.
func (a *Aggregator) History() IndicatorHistory {
	return IndicatorHistory{
		RSI:    copySlice(a.rsiHistory),
		SMA:    copySlice(a.smaHistory),
		EMA:    copySlice(a.emaHistory),
		Prices: copySlice(a.priceHistory),
	}
}

func appendWithLimit(slice []float64, val float64, limit int) []float64 {
	slice = append(slice, val)
	if len(slice) > limit {
		slice = slice[len(slice)-limit:]
	}
	return slice
}

func copySlice(src []float64) []float64 {
	if src == nil {
		return nil
	}
	dst := make([]float64, len(src))
	copy(dst, src)
	return dst
}

func maxInt(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
