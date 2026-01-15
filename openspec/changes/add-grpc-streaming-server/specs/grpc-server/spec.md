# Spec: gRPC Server

## ADDED Requirements

### Requirement: Server Initialization
The system SHALL start a gRPC server on port 50051 when executing `cmd/server/main.go`.

#### Scenario: Server Startup
- **WHEN** the server process is initiated
- **THEN** a gRPC server SHALL bind to `:50051`
- **AND** the server SHALL begin accepting incoming connections
- **AND** the server SHALL log "gRPC server listening on :50051"

### Requirement: StreamPrices RPC Method
The system SHALL implement a bidirectional streaming RPC method `StreamPrices` in `internal/grpc/server/handler.go`.

#### Scenario: StreamPrices Signature
- **GIVEN** a gRPC client connects to the server
- **WHEN** the client calls `StreamPrices`
- **THEN** the server SHALL accept a stream of `IndicatorRequest` messages from the client
- **AND** the server SHALL send a stream of `PriceUpdate` and `IndicatorUpdate` messages to the client

#### Scenario: Initial Price Subscription
- **GIVEN** a client initiates `StreamPrices`
- **THEN** the server SHALL immediately begin streaming `PriceUpdate` messages for subscribed symbols
- **AND** the first `PriceUpdate` SHALL contain current market data for BTCUSDT

### Requirement: Protocol Buffer Definitions
The system SHALL define the following messages in the Protocol Buffer schema.

#### Scenario: IndicatorRequest Message
- **GIVEN** the proto schema requires client configuration capability
- **WHEN** defining the `IndicatorRequest` message
- **THEN** it SHALL include:
  - `symbol` (string): Trading pair symbol (e.g., "btcusdt")
  - `indicator_type` (enum): RSI, SMA, or EMA
  - `period` (int32): Calculation period (default: 14 for RSI)
  - `action` (enum): SUBSCRIBE or UNSUBSCRIBE

#### Scenario: PriceUpdate Message
- **GIVEN** the proto schema requires Binance ticker data
- **WHEN** defining the `PriceUpdate` message
- **THEN** it SHALL include:
  - `symbol` (string): Trading pair symbol
  - `price` (string): Current price
  - `volume` (string): 24h trading volume
  - `timestamp` (int64): Unix timestamp in milliseconds

#### Scenario: IndicatorUpdate Message
- **GIVEN** the proto schema requires calculated indicator delivery
- **WHEN** defining the `IndicatorUpdate` message
- **THEN** it SHALL include:
  - `symbol` (string): Trading pair symbol
  - `indicator_type` (enum): RSI, SMA, or EMA
  - `value` (double): Calculated indicator value
  - `period` (int32): Calculation period used
  - `timestamp` (int64): Unix timestamp in milliseconds

### Requirement: Concurrent Client Handling
The system SHALL support multiple concurrent gRPC client connections.

#### Scenario: Multiple Client Connections
- **GIVEN** the gRPC server is running
- **WHEN** multiple clients connect to `StreamPrices`
- **THEN** each client SHALL receive independent streams
- **AND** price updates SHALL be fanned out to all connected clients
- **AND** indicator calculations SHALL be independent per client configuration

### Requirement: Connection Lifecycle
The system SHALL properly manage gRPC stream lifecycle.

#### Scenario: Client Disconnection
- **WHEN** a gRPC client disconnects from `StreamPrices`
- **THEN** the server SHALL clean up the client's stream registration
- **AND** the server SHALL stop sending updates to that client
- **AND** the server SHALL NOT crash or interrupt other clients

#### Scenario: Client Sends Invalid Request
- **WHEN** a client sends an `IndicatorRequest` with invalid data
- **THEN** the server SHALL log a warning
- **AND** the server SHALL continue the stream without interruption
- **AND** the invalid request SHALL be ignored
