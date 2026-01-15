# Spec: Binance Integration

## ADDED Requirements

### Requirement: Binance WebSocket Connection
The system SHALL establish and maintain a WebSocket connection to Binance for real-time market data.

#### Scenario: Initial Connection
- **WHEN** the gRPC server starts
- **THEN** the system SHALL connect to `wss://stream.binance.com:9443/stream`
- **AND** the system SHALL subscribe to the `btcusdt@ticker` stream
- **AND** the connection SHALL be logged as "Connected to Binance WebSocket"

#### Scenario: Connection Parameters
- **GIVEN** the Binance WebSocket endpoint requirements
- **WHEN** establishing the connection
- **THEN** the system SHALL use the combined stream endpoint `/stream`
- **AND** the system SHALL subscribe to streams using the `streams` query parameter
- **AND** the system SHALL request `btcusdt@ticker` for initial price data

### Requirement: Price Data Reception
The system SHALL receive and parse price update messages from Binance.

#### Scenario: Ticker Message Reception
- **WHEN** Binance sends a ticker message
- **THEN** the system SHALL parse the JSON message
- **AND** the system SHALL extract: symbol, price, volume, timestamp
- **AND** the parsed data SHALL be converted to a `PriceUpdate` message

#### Scenario: Price Update Rate
- **GIVEN** the Binance ticker stream rate
- **WHEN** price updates arrive from Binance
- **THEN** the system SHALL process updates at their native rate (typically 100ms intervals)
- **AND** the system SHALL NOT drop updates under normal load

### Requirement: Price Data Distribution
The system SHALL distribute received price updates to all connected gRPC clients.

#### Scenario: Fan-Out Distribution
- **WHEN** a `PriceUpdate` is received from Binance
- **THEN** the system SHALL distribute it to all active `StreamPrices` streams
- **AND** each client SHALL receive the update in real-time
- **AND** the distribution SHALL use non-blocking sends to prevent backpressure

### Requirement: Reconnection Handling
The system SHALL handle Binance WebSocket disconnections gracefully.

#### Scenario: Unexpected Disconnection
- **WHEN** the Binance WebSocket connection is lost
- **THEN** the system SHALL log "Binance WebSocket disconnected"
- **AND** the system SHALL attempt reconnection with exponential backoff
- **AND** the initial backoff delay SHALL be 1 second
- **AND** the maximum backoff delay SHALL be 30 seconds

#### Scenario: Reconnection Success
- **WHEN** the reconnection attempt succeeds
- **THEN** the system SHALL log "Reconnected to Binance WebSocket"
- **AND** the system SHALL resubscribe to price streams
- **AND** the system SHALL resume distribution to all connected clients

#### Scenario: Maximum Retries Exceeded
- **WHEN** reconnection attempts exceed 5 consecutive failures
- **THEN** the system SHALL log "Binance WebSocket connection failed after 5 attempts"
- **AND** the system SHALL enter a degraded state
- **AND** the system SHALL continue to accept gRPC connections but send no updates

### Requirement: Connection Health Monitoring
The system SHALL monitor the health of the Binance WebSocket connection.

#### Scenario: Heartbeat Check
- **GIVEN** the Binance WebSocket connection is established
- **WHEN** no message is received for 60 seconds
- **THEN** the system SHALL consider the connection stale
- **AND** the system SHALL attempt to reconnect

#### Scenario: Ping/Pong Handling
- **WHEN** the system sends a ping to Binance
- **THEN** the system SHALL expect a pong within 10 seconds
- **AND** the absence of a pong SHALL trigger reconnection
