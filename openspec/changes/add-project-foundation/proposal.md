# Change: Initialize Project Foundation

## Why
Set up the foundational Go project structure, gRPC service definitions, and Binance WebSocket client to enable real-time cryptocurrency market data streaming as the core infrastructure for the QuantaCode trading analysis CLI.

## What Changes
- Initialize Go module with basic project structure (cmd/server, cmd/cli, internal/domain, internal/infra)
- Add Go-specific .gitignore and project README
- Define Protocol Buffer schema for market data (PriceUpdate message with symbol, price, volume, timestamp)
- Generate Go code from proto definitions with protoc
- Implement gRPC service definition for streaming price updates (StreamPrices)
- Create Binance WebSocket client connecting to wss://stream.binance.com:9443/ws/btcusdt@ticker
- Implement JSON ticker parsing and PriceUpdate struct mapping
- Add reconnection logic for WebSocket connection failures

## Impact
- Affected specs: grpc-market-data, binance-client
- New directories: cmd/, internal/, proto/
- New capabilities: gRPC streaming service, Binance market data ingestion
