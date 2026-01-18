## 1. Indicator Infrastructure
- [x] 1.1 Create circular buffer implementation for price history
- [x] 1.2 Define common indicator interface (if needed for consistency)

## 2. RSI Implementation
- [x] 2.1 Create internal/domain/indicators/rsi.go file
- [x] 2.2 Implement RSI calculation with period parameter
- [x] 2.3 Integrate circular buffer for price history
- [x] 2.4 Add table-driven unit tests with known input/output values

## 3. SMA Implementation
- [x] 3.1 Create internal/domain/indicators/sma.go file
- [x] 3.2 Implement SMA calculation with period parameter
- [x] 3.3 Use circular buffer for price history
- [x] 3.4 Add table-driven unit tests with known input/output values

## 4. EMA Implementation
- [x] 4.1 Create internal/domain/indicators/ema.go file
- [x] 4.2 Implement EMA calculation with period parameter
- [x] 4.3 Use circular buffer for price history
- [x] 4.4 Add table-driven unit tests with known input/output values

## 5. Indicator Aggregator
- [x] 5.1 Create internal/domain/indicators/aggregator.go file
- [x] 5.2 Implement price buffer management
- [x] 5.3 Implement update-all-indicators logic on new price
- [x] 5.4 Define aggregated values return struct
- [x] 5.5 Add unit tests for aggregator

## 6. Validation
- [x] 6.1 Run all unit tests with `go test ./internal/domain/indicators/...`
- [x] 6.2 Verify RSI calculations match expected values (e.g., standard 14-period RSI)
- [x] 6.3 Verify SMA/EMA calculations with manual calculations
- [x] 6.4 Test edge cases (empty buffer, single price, insufficient data)
- [x] 6.5 Run `go vet` and static analysis
