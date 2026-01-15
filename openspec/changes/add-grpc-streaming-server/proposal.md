# Change: Implement gRPC Streaming Server with Binance Integration

## Why
Implement the gRPC server infrastructure to serve real-time cryptocurrency price data and technical indicators to CLI clients. This establishes the server-side foundation for the bidirectional streaming architecture, enabling clients to receive live price updates from Binance while simultaneously streaming calculated indicators (RSI, SMA, EMA) based on their subscription preferences. The graceful shutdown mechanism ensures clean termination during deployments or unexpected interruptions.

## What Changes

### Core Server Implementation
- Create `cmd/server/main.go` with gRPC server initialization on port 50051
- Implement `internal/grpc/server/handler.go` with `StreamPrices` bidirectional streaming method
- Integrate Binance WebSocket client for real-time price data subscription
- Add signal handling for SIGTERM/SIGINT with graceful shutdown (5s timeout)

### Protocol Buffer Extensions
- Add `IndicatorRequest` message to support client-configurable indicator subscriptions
- Extend streaming to support simultaneous `PriceUpdate` and calculated indicator delivery
- Define indicator messages (RSI, SMA, EMA) with value and metadata fields

### Indicator Streaming
- Implement real-time RSI calculation and streaming
- Implement real-time SMA calculation and streaming  
- Implement real-time EMA calculation and streaming
- Add client-configurable indicator parameters (period, symbol)

### Operational Reliability
- Implement graceful shutdown: close Binance connection, drain active streams, shutdown gRPC server
- Add connection health monitoring and reconnection logic
- Configure server-side stream management for multiple concurrent clients

## Impact
- Affected specs: grpc-market-data, binance-client, technical-indicators
- New files: `cmd/server/main.go`, `internal/grpc/server/handler.go`, `proto/*.proto` extensions
- New capabilities: gRPC bidirectional streaming, indicator calculation service, graceful shutdown
- Dependencies: Extends existing Binance client and indicator calculations from project foundation
