## ADDED Requirements

### Requirement: AIAnalysis Proto Message
The system SHALL define AIAnalysis message in the Protocol Buffer schema containing symbol, analysis text, timestamp, and streaming status.

#### Scenario: Define AIAnalysis message
- **WHEN** proto file defines AIAnalysis
- **THEN** it contains: string symbol, string analysis, int64 timestamp, bool is_streaming

### Requirement: StreamAIAnalysis gRPC Method
The system SHALL provide a StreamAIAnalysis gRPC service method that returns a stream of AIAnalysis messages.

#### Scenario: Stream AI responses
- **WHEN** client calls StreamAIAnalysis
- **THEN** it receives stream of AIAnalysis messages with analysis text

### Requirement: Price Change Trigger
The system SHALL trigger AI analysis only when price changes by more than 1% from previous value.

#### Scenario: Price exceeds 1% threshold
- **WHEN** current price differs from previous price by >1%
- **THEN** AI analysis is initiated

#### Scenario: Price below 1% threshold
- **WHEN** current price differs from previous price by <=1%
- **THEN** no AI analysis is triggered

### Requirement: OpenRouter Integration
The system SHALL call OpenRouter API to generate AI analysis when triggered by price change.

#### Scenario: Request AI analysis
- **WHEN** price change triggers analysis
- **THEN** server calls OpenRouter with formatted analysis context

### Requirement: Response Streaming
The system SHALL stream AI analysis responses to the gRPC client in real-time.

#### Scenario: Stream partial responses
- **WHEN** OpenRouter returns partial response
- **THEN** server sends AIAnalysis with is_streaming=true

#### Scenario: Stream final response
- **WHEN** OpenRouter completes response
- **THEN** server sends final AIAnalysis with is_streaming=false

### Requirement: Financial Disclaimer
The system SHALL append "Not financial advice" disclaimer to all AI analysis responses.

#### Scenario: Add disclaimer to response
- **WHEN** sending AIAnalysis to client
- **THEN** analysis text includes disclaimer statement
