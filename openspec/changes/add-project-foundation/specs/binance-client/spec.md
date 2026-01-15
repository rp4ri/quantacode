## ADDED Requirements

### Requirement: Binance WebSocket Connection
The system SHALL connect to Binance WebSocket endpoint wss://stream.binance.com:9443/ws/btcusdt@ticker to receive real-time market data.

#### Scenario: Successful connection
- **WHEN** the client initiates connection to Binance WebSocket
- **THEN** it establishes a connection and receives ticker data

### Requirement: Ticker Message Parsing
The system SHALL parse JSON ticker messages from Binance and map them to the PriceUpdate struct.

#### Scenario: Parse BTCUSDT ticker
- **WHEN** a JSON ticker message is received from Binance
- **THEN** it is parsed into PriceUpdate with symbol, price, volume, and timestamp fields

### Requirement: Automatic Reconnection
The system SHALL automatically reconnect to the Binance WebSocket when connection errors occur.

#### Scenario: Reconnect on disconnect
- **WHEN** the WebSocket connection is lost
- **THEN** the client automatically reconnects with exponential backoff retry logic

### Requirement: Error Handling
The system SHALL handle and log WebSocket errors without crashing the application.

#### Scenario: Log connection error
- **WHEN** a WebSocket error occurs
- **THEN** the error is logged and reconnection is attempted
