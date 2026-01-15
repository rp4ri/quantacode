## ADDED Requirements

### Requirement: Analysis Context Builder
The system SHALL provide a component that formats price data and technical indicators into a structured prompt for AI analysis.

#### Scenario: Build analysis context
- **WHEN** context builder receives price and indicators
- **THEN** it creates structured prompt with symbol, price, RSI, SMA, EMA values

### Requirement: Include Current Values
The system SHALL include current price action, technical indicator values, and trends in the AI prompt.

#### Scenario: Format current market state
- **WHEN** building analysis context
- **THEN** prompt contains current price, RSI value (overbought/oversold status), SMA/EMA trend direction

### Requirement: Explain Why Instruction
The system SHALL include explicit instruction in the prompt for the AI to explain "why" not just "what".

#### Scenario: Add explanation directive
- **WHEN** building analysis context
- **THEN** prompt instructs AI to explain reasoning behind market conditions

### Requirement: Trend Analysis
The system SHALL analyze price trends (up, down, sideways) and include trend direction in the prompt.

#### Scenario: Determine trend direction
- **WHEN** comparing current price to SMA
- **THEN** trend is classified as up (price > SMA), down (price < SMA), or sideways (price ≈ SMA)

### Requirement: Context Structuring
The system SHALL structure the analysis context in a clear, machine-readable format for the AI.

#### Scenario: Format structured prompt
- **WHEN** building analysis context
- **THEN** prompt uses sections for "Current State", "Technical Indicators", "Trend Analysis", and "Analysis Request"
