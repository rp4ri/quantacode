## ADDED Requirements

### Requirement: RSI Calculation
The system SHALL implement Relative Strength Index (RSI) calculation with a configurable period parameter using a circular buffer for price history.

#### Scenario: Calculate 14-period RSI
- **WHEN** 15 price updates are provided with period=14
- **THEN** RSI value is calculated and returned between 0-100 range

#### Scenario: Handle insufficient data
- **WHEN** fewer than period+1 prices are available
- **THEN** RSI returns 0 indicating not ready

### Requirement: SMA Calculation
The system SHALL implement Simple Moving Average (SMA) calculation with a configurable period parameter using a circular buffer for price history.

#### Scenario: Calculate 20-period SMA
- **WHEN** 20 price updates are provided with period=20
- **THEN** SMA equals arithmetic mean of the last 20 prices

#### Scenario: Handle insufficient data
- **WHEN** fewer than period prices are available
- **THEN** SMA returns 0 indicating not ready

### Requirement: EMA Calculation
The system SHALL implement Exponential Moving Average (EMA) calculation with a configurable period parameter using a circular buffer for price history.

#### Scenario: Calculate 12-period EMA
- **WHEN** multiple price updates are provided with period=12
- **THEN** EMA applies smoothing factor 2/(period+1) to new prices

#### Scenario: Initialize EMA with SMA
- **WHEN** EMA receives first period prices
- **THEN** initial EMA value equals SMA of first period prices

### Requirement: Indicator Aggregation
The system SHALL provide an indicator aggregator that maintains a price buffer and updates all registered indicators on each new price.

#### Scenario: Update all indicators on price
- **WHEN** aggregator receives a new price update
- **THEN** RSI, SMA, and EMA are all updated and available

#### Scenario: Return aggregated values
- **WHEN** all indicators have sufficient data
- **THEN** aggregator returns struct with current RSI, SMA, and EMA values

### Requirement: Indicator Testing
The system SHALL include comprehensive unit tests using table-driven test patterns for all indicator calculations.

#### Scenario: RSI test with known values
- **WHEN** RSI is tested with predefined price sequence
- **THEN** calculated RSI matches expected value within tolerance

#### Scenario: SMA test with known values
- **WHEN** SMA is tested with predefined price sequence
- **THEN** calculated SMA matches expected value exactly

#### Scenario: EMA test with known values
- **WHEN** EMA is tested with predefined price sequence
- **THEN** calculated EMA matches expected value within tolerance
