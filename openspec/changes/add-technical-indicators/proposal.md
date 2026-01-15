# Change: Add Technical Indicators

## Why
Implement core technical analysis indicators (RSI, SMA, EMA) with efficient price history management using circular buffers, enabling real-time market analysis for the QuantaCode trading CLI.

## What Changes
- Implement RSI (Relative Strength Index) calculation with configurable period parameter
- Implement SMA (Simple Moving Average) calculation with configurable period parameter
- Implement EMA (Exponential Moving Average) calculation with configurable period parameter
- Use circular buffer pattern for efficient price history management
- Create indicator aggregator that maintains price buffer and updates all indicators on new price
- Return aggregated indicator values struct
- Add comprehensive unit tests with table-driven tests for all indicators

## Impact
- Affected specs: indicators
- New files: internal/domain/indicators/rsi.go, sma.go, ema.go, aggregator.go
- New capability: Real-time technical indicator calculations
