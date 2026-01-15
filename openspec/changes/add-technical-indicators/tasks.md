## 1. Indicator Infrastructure
- [ ] 1.1 Create circular buffer implementation for price history
- [ ] 1.2 Define common indicator interface (if needed for consistency)

## 2. RSI Implementation
- [ ] 2.1 Create internal/domain/indicators/rsi.go file
- [ ] 2.2 Implement RSI calculation with period parameter
- [ ] 2.3 Integrate circular buffer for price history
- [ ] 2.4 Add table-driven unit tests with known input/output values

## 3. SMA Implementation
- [ ] 3.1 Create internal/domain/indicators/sma.go file
- [ ] 3.2 Implement SMA calculation with period parameter
- [ ] 3.3 Use circular buffer for price history
- [ ] 3.4 Add table-driven unit tests with known input/output values

## 4. EMA Implementation
- [ ] 4.1 Create internal/domain/indicators/ema.go file
- [ ] 4.2 Implement EMA calculation with period parameter
- [ ] 4.3 Use circular buffer for price history
- [ ] 4.4 Add table-driven unit tests with known input/output values

## 5. Indicator Aggregator
- [ ] 5.1 Create internal/domain/indicators/aggregator.go file
- [ ] 5.2 Implement price buffer management
- [ ] 5.3 Implement update-all-indicators logic on new price
- [ ] 5.4 Define aggregated values return struct
- [ ] 5.5 Add unit tests for aggregator

## 6. Validation
- [ ] 6.1 Run all unit tests with `go test ./internal/domain/indicators/...`
- [ ] 6.2 Verify RSI calculations match expected values (e.g., standard 14-period RSI)
- [ ] 6.3 Verify SMA/EMA calculations with manual calculations
- [ ] 6.4 Test edge cases (empty buffer, single price, insufficient data)
- [ ] 6.5 Run `go vet` and static analysis
