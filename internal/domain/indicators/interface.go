package indicators

// Indicator defines a common API for technical indicators.
type Indicator interface {
	Update(price float64) float64
	Value() float64
	Period() int
}
