## ADDED Requirements

### Requirement: Protocol Buffer Schema
The system SHALL define a Protocol Buffer schema for market data messages that includes symbol, price, volume, and timestamp fields.

#### Scenario: PriceUpdate message definition
- **WHEN** the proto file defines PriceUpdate
- **THEN** it contains: string symbol, double price, double volume, int64 timestamp

### Requirement: Go Code Generation
The system SHALL generate Go code from Protocol Buffer definitions using protoc and protoc-gen-go tools.

#### Scenario: Generated code structure
- **WHEN** protoc is run on market_data.proto
- **THEN** Go files are generated in the proto/ directory with message types and service interfaces

### Requirement: gRPC Streaming Service
The system SHALL provide a gRPC service definition with a StreamPrices method that supports bidirectional streaming of PriceUpdate messages.

#### Scenario: StreamPrices service method
- **WHEN** a client calls StreamPrices
- **THEN** it returns a stream of PriceUpdate messages from the server
