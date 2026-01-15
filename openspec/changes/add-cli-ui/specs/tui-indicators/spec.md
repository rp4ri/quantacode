## ADDED Requirements

### Requirement: Indicators Panel Layout
The system SHALL provide a right sidebar panel displaying technical indicators.

#### Scenario: Display indicators panel
- **WHEN** TUI is running
- **THEN** right sidebar shows RSI, SMA, EMA values with labels

### Requirement: RSI Display
The system SHALL display RSI (Relative Strength Index) value in the indicators panel.

#### Scenario: Show RSI value
- **WHEN** RSI calculation is available
- **THEN** panel displays "RSI: <value>"

### Requirement: SMA Display
The system SHALL display SMA (Simple Moving Average) value in the indicators panel.

#### Scenario: Show SMA value
- **WHEN** SMA calculation is available
- **THEN** panel displays "SMA14: <value>"

### Requirement: EMA Display
The system SHALL display EMA (Exponential Moving Average) value in the indicators panel.

#### Scenario: Show EMA value
- **WHEN** EMA calculation is available
- **THEN** panel displays "EMA12: <value>"

### Requirement: RSI Color Coding
The system SHALL color-code RSI values: red for overbought (>70), blue for oversold (<30), white for neutral.

#### Scenario: Overbought RSI
- **WHEN** RSI value > 70
- **THEN** RSI displays in red

#### Scenario: Oversold RSI
- **WHEN** RSI value < 30
- **THEN** RSI displays in blue

#### Scenario: Neutral RSI
- **WHEN** RSI value is between 30 and 70
- **THEN** RSI displays in white

### Requirement: Real-time Updates
The system SHALL update indicator values in real-time as new price data arrives.

#### Scenario: Update on price change
- **WHEN** new price update triggers indicator recalculation
- **THEN** indicator panel displays updated values

#### Scenario: Maintain previous values
- **WHEN** indicator is not ready (insufficient data)
- **THEN** panel shows 0 or previous value

### Requirement: Indicator Formatting
The system SHALL format indicator values with appropriate precision for readability.

#### Scenario: Format RSI values
- **WHEN** RSI is displayed
- **THEN** it shows 1 decimal place (e.g., "75.2")

#### Scenario: Format SMA/EMA values
- **WHEN** SMA or EMA is displayed
- **THEN** they show 2 decimal places (e.g., "$50,123.45")

### Requirement: Panel Styling
The system SHALL style the indicators panel with Lipgloss for visual separation.

#### Scenario: Apply panel border
- **WHEN** panel is rendered
- **THEN** it has a border separating from main chat area

#### Scenario: Apply padding
- **WHEN** panel is rendered
- **THEN** indicators have padding for readability
