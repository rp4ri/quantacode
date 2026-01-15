# Spec: Indicator Streaming

## ADDED Requirements

### Requirement: Indicator Request Handling
The system SHALL process `IndicatorRequest` messages from clients to configure indicator subscriptions.

#### Scenario: Subscribe to RSI
- **WHEN** a client sends an `IndicatorRequest` with:
  - `symbol`: "btcusdt"
  - `indicator_type`: RSI
  - `period`: 14
  - `action`: SUBSCRIBE
- **THEN** the server SHALL begin calculating RSI for BTCUSDT with period 14
- **AND** the server SHALL stream `IndicatorUpdate` messages to the client

#### Scenario: Subscribe to SMA
- **WHEN** a client sends an `IndicatorRequest` with:
  - `symbol`: "btcusdt"
  - `indicator_type`: SMA
  - `period`: 20
  - `action`: SUBSCRIBE
- **THEN** the server SHALL begin calculating SMA for BTCUSDT with period 20
- **AND** the server SHALL stream `IndicatorUpdate` messages to the client

#### Scenario: Subscribe to EMA
- **WHEN** a client sends an `IndicatorRequest` with:
  - `symbol`: "ethusdt"
  - `indicator_type`: EMA
  - `period`: 12
  - `action`: SUBSCRIBE
- **THEN** the server SHALL begin calculating EMA for ETHUSDT with period 12
- **AND** the server SHALL stream `IndicatorUpdate` messages to the client

#### Scenario: Unsubscribe from Indicator
- **WHEN** a client sends an `IndicatorRequest` with `action`: UNSUBSCRIBE
- **THEN** the server SHALL stop calculating that indicator for the client
- **AND** the server SHALL remove the indicator subscription
- **AND** no further `IndicatorUpdate` messages for that indicator SHALL be sent

### Requirement: RSI Calculation
The system SHALL calculate and stream Relative Strength Index (RSI) values.

#### Scenario: RSI Calculation Algorithm
- **GIVEN** a price history buffer for a symbol
- **WHEN** calculating RSI
- **THEN** the system SHALL use the standard 14-period RSI formula
- **AND** the system SHALL calculate average gains and average losses over the period
- **AND** the system SHALL compute RS = average_gain / average_loss
- **AND** the system SHALL compute RSI = 100 - (100 / (1 + RS))

#### Scenario: RSI Initial Calculation
- **WHEN** fewer than 14 price updates are available
- **THEN** the system SHALL NOT send RSI values until 14 data points exist
- **AND** the system SHALL accumulate price data silently
- **AND** once 14 points exist, RSI SHALL be calculated and streamed

#### Scenario: RSI Update Frequency
- **WHEN** a new price update arrives for a subscribed symbol
- **THEN** the RSI SHALL be recalculated
- **AND** the new RSI value SHALL be wrapped in an `IndicatorUpdate`
- **AND** the update SHALL be streamed to subscribed clients within 100ms

### Requirement: SMA Calculation
The system SHALL calculate and stream Simple Moving Average (SMA) values.

#### Scenario: SMA Calculation Algorithm
- **GIVEN** a price history buffer for a symbol
- **WHEN** calculating SMA
- **THEN** the system SHALL sum the last N closing prices
- **AND** the system SHALL divide by N (the period)
- **AND** the result SHALL be the SMA value

#### Scenario: SMA Initial Calculation
- **WHEN** fewer than N (period) price updates are available
- **THEN** the system SHALL NOT send SMA values until N data points exist
- **AND** the system SHALL accumulate price data silently
- **AND** once N points exist, SMA SHALL be calculated and streamed

#### Scenario: SMA Rolling Window
- **WHEN** a new price update arrives for a subscribed symbol
- **THEN** the oldest price in the window SHALL be removed
- **AND** the new price SHALL be added to the window
- **AND** the SMA SHALL be recalculated
- **AND** the new SMA value SHALL be streamed to subscribed clients

### Requirement: EMA Calculation
The system SHALL calculate and stream Exponential Moving Average (EMA) values.

#### Scenario: EMA Calculation Algorithm
- **GIVEN** a price history buffer for a symbol
- **WHEN** calculating EMA
- **THEN** the system SHALL use the smoothing factor: α = 2 / (N + 1)
- **AND** the system SHALL compute: EMA_t = α × Price_t + (1 - α) × EMA_{t-1}
- **AND** the first EMA SHALL equal the first SMA (or first price if unavailable)

#### Scenario: EMA Initial Calculation
- **WHEN** fewer than 2 price updates are available
- **THEN** the system SHALL NOT send EMA values until sufficient data exists
- **AND** the system SHALL accumulate price data silently
- **AND** once 2+ points exist, EMA SHALL be calculated and streamed

#### Scenario: EMA Update Frequency
- **WHEN** a new price update arrives for a subscribed symbol
- **THEN** the EMA SHALL be recalculated using the new price
- **AND** the new EMA value SHALL be streamed to subscribed clients

### Requirement: Indicator Error Handling
The system SHALL handle indicator calculation errors gracefully.

#### Scenario: Insufficient Data
- **WHEN** an indicator cannot be calculated due to insufficient data
- **THEN** the system SHALL log a debug message
- **AND** the system SHALL skip that indicator update
- **AND** the system SHALL continue processing other indicators
- **AND** the gRPC stream SHALL NOT be interrupted

#### Scenario: Invalid Period
- **WHEN** a client requests an indicator with period <= 0
- **THEN** the system SHALL use a default period (14 for RSI, 20 for SMA/EMA)
- **AND** the system SHALL log a warning
- **AND** the subscription SHALL proceed with default period

#### Scenario: Price Data Gap
- **WHEN** price updates are missed or delayed
- **THEN** the system SHALL continue calculating with available data
- **AND** the system SHALL NOT crash on missing data points
- **AND** the system SHALL log a warning if data gaps exceed expected intervals

### Requirement: Memory Management
The system SHALL manage memory efficiently for price history buffers.

#### Scenario: Price History Limit
- **GIVEN** the maximum period across all active indicator subscriptions
- **WHEN** storing price history
- **THEN** the system SHALL maintain at least that many price points per symbol
- **AND** the system SHALL NOT store more than 200 price points per symbol
- **AND** older prices beyond the maximum period SHALL be discarded

#### Scenario: Per-Symbol Buffers
- **GIVEN** multiple symbols being tracked
- **WHEN** storing price history
- **THEN** each symbol SHALL have its own independent buffer
- **AND** buffers SHALL be created on first subscription to a symbol
- **AND** buffers SHALL be cleaned up when no clients are subscribed
