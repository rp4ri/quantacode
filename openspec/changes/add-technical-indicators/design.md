## Context
Technical indicators require maintaining price history for calculations (e.g., 14-period RSI needs 14 price points). A naive approach using slices would cause O(n) allocations per update, which is inefficient for real-time streaming with frequent updates.

## Goals / Non-Goals
- Goals: Efficient O(1) price history updates, reusable indicator calculations, aggregated indicator access
- Non-Goals: Supporting multiple symbols in single aggregator, time-weighted indicators, advanced indicators (MACD, Bollinger Bands) in this phase

## Decisions

### Circular Buffer for Price History
- **Decision**: Use fixed-size circular buffer to store price history
- **Rationale**: O(1) write operations, no memory allocations per update, minimal memory overhead
- **Implementation**: Pre-allocated slice with write pointer that wraps around

### Indicator Calculation Pattern
- **Decision**: Each indicator maintains its own circular buffer sized to its period
- **Rationale**: Simpler implementation, each indicator self-contained, easier to test independently
- **Alternative**: Shared buffer in aggregator - rejected due to complexity and coupling

### Aggregator Component
- **Decision**: Create aggregator that maintains references to indicators and updates all on price update
- **Rationale**: Single point for price updates, returns all indicator values in struct
- **Alternative**: Client manages indicators directly - rejected for better encapsulation

### Period Parameter Design
- **Decision**: Period passed at indicator construction, immutable after creation
- **Rationale**: Period affects buffer size, should be fixed at initialization

## Data Structures

```go
// CircularBuffer implements O(1) push/pop operations
type CircularBuffer struct {
    buffer []float64
    size   int
    index  int
}

// Indicator interface for consistent API
type Indicator interface {
    Update(price float64) float64
    Value() float64
    Period() int
}

// AggregatedIndicatorValues struct
type AggregatedIndicatorValues struct {
    RSI    float64
    SMA14  float64
    EMA12  float64
    // ... more indicators as needed
}
```

## Risks / Trade-offs

### Memory Usage
- **Risk**: Multiple circular buffers (one per indicator) duplicate price data
- **Mitigation**: Memory overhead is minimal (<100KB for 3 indicators with 100-period history)
- **Trade-off**: Accepting slight memory duplication for code simplicity and testability

### Insufficient Data Handling
- **Risk**: Indicators return invalid values when buffer not full (e.g., RSI needs at least period+1 prices)
- **Mitigation**: Return 0 or NaN for insufficient data, document clearly in API
- **Trade-off**: Prefer explicit "not ready" state over returning partial calculations

### Floating Point Precision
- **Risk**: Cumulative errors in EMA calculation over time
- **Mitigation**: Use float64 (not float32), validate with known test cases
- **Trade-off**: Accept IEEE 754 precision limits, sufficient for trading analysis

## Migration Plan
N/A - new capability, no existing code to migrate

## Open Questions
- Should indicator values return 0 or NaN when insufficient data?
- Should aggregator initialize with default indicator periods (RSI=14, SMA=20, EMA=12)?
- Should indicators be thread-safe? (assuming single-threaded for now)
