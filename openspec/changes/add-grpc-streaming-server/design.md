# Design: gRPC Streaming Server Architecture

## Overview
This document captures architectural decisions for the gRPC streaming server implementation with Binance integration, bidirectional indicator streaming, and graceful shutdown handling.

## Architectural Decisions

### 1. Bidirectional Streaming Pattern

**Decision**: Use gRPC bidirectional streaming for the `StreamPrices` RPC method.

**Rationale**:
- Clients can send `IndicatorRequest` messages at any time to modify indicator subscriptions without reconnecting
- Server can interleave `PriceUpdate` and calculated `IndicatorUpdate` messages on the same stream
- Reduces connection overhead and enables real-time configuration changes

**Trade-offs**:
- Increases implementation complexity compared to simple server-side streaming
- Requires careful synchronization when multiple concurrent indicator calculations share price data

**Alternative Considered**: Separate unary requests for indicator configuration with server-side streaming for data. Rejected due to latency in configuration updates and increased connection management.

### 2. Binance Connection Management

**Decision**: Single shared Binance WebSocket connection per server instance, multiplexed to all connected gRPC clients.

**Rationale**:
- Binance limits: 5 concurrent WebSocket connections, 1200 requests/minute
- Shared connection reduces resource usage and respects rate limits
- All clients receive same symbol data; no per-client duplication

**Trade-offs**:
- Single point of failure for Binance connection
- All clients receive identical symbol subscriptions
- Requires fan-out mechanism for price updates to multiple gRPC streams

### 3. Indicator Calculation Strategy

**Decision**: Server-side indicator calculations using sliding window of recent price data.

**Rationale**:
- Keeps indicator logic centralized and consistent across all clients
- Reduces client-side computation requirements
- Enables future extensibility (e.g., server-side alerts)

**Trade-offs**:
- Server CPU usage increases with number of symbols and indicators
- Requires careful memory management for price history buffers
- Calculation latency may introduce lag vs. client-side computation

**Implementation**:
- Price history buffer per symbol (configurable window size)
- Incremental calculation updates as new prices arrive
- Batch indicator updates to reduce gRPC message frequency (e.g., every 100ms)

### 4. Graceful Shutdown Sequence

**Decision**: Three-phase shutdown with configurable timeout (5s default).

**Phase 1 - Signal Receipt**:
- Catch SIGTERM/SIGINT via context cancellation
- Set shutdown flag to prevent new stream accepts
- Allow existing streams to continue during drain window

**Phase 2 - Stream Drain**:
- Notify all active streams of impending shutdown
- Continue processing for drain timeout (5s)
- Close Binance WebSocket connection
- Complete final indicator calculations

**Phase 3 - Server Stop**:
- Force close any remaining streams after timeout
- Graceful gRPC server shutdown with `GracefulStop()`
- Release all resources

**Rationale**:
- Ensures in-flight messages complete before termination
- Prevents data loss during rolling deployments
- Avoids hard cutoffs that could corrupt indicator calculations

### 5. Stream Management

**Decision**: Per-client stream state tracking with shared price data distributor.

**Components**:
- `StreamRegistry`: Tracks active gRPC streams per client
- `PriceDistributor`: Fans out Binance price updates to all registered streams
- `IndicatorEngine`: Calculates indicators based on price updates and subscription state

**Concurrency Model**:
- Channel-based communication between components
- Mutex protection for shared state (stream registry, price history)
- Goroutines per gRPC stream for bidirectional message handling

## Component Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      cmd/server/main.go                          │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │  gRPC Server (:50051)                                       ││
│  │  ┌─────────────────┐  ┌──────────────┐  ┌─────────────────┐ ││
│  │  │ Signal Handler  │  │ StreamRegistry│  │ Shutdown Coord  │ ││
│  │  │ (SIGTERM/SIGINT)│  │              │  │                 │ ││
│  │  └────────┬────────┘  └──────┬───────┘  └────────┬────────┘ ││
│  │           │                  │                   │          ││
│  │           ▼                  ▼                   ▼          ││
│  │  ┌─────────────────────────────────────────────────────┐   ││
│  │  │              Binance WebSocket Client               │   ││
│  │  │  wss://stream.binance.com:9443/stream               │   ││
│  │  └────────────────────────┬────────────────────────────┘   ││
│  │                           │                                 ││
│  │                           ▼                                 ││
│  │  ┌─────────────────────────────────────────────────────┐   ││
│  │  │              Price Distributor                       │   ││
│  │  │         (Fan-out to all active streams)             │   ││
│  │  └────────────────────────┬────────────────────────────┘   ││
│  │                           │                                 ││
│  │              ┌────────────┼────────────┐                   ││
│  │              ▼            ▼            ▼                   ││
│  │  ┌─────────────────┐ ┌─────────────────┐ ┌──────────────┐ ││
│  │  │ RSI Calculator  │ │ SMA Calculator  │ │ EMA Calculator│ ││
│  │  └────────┬────────┘ └────────┬────────┘ └──────┬───────┘ ││
│  │           │                  │                  │          ││
│  │           └──────────────────┼──────────────────┘          ││
│  │                              ▼                             ││
│  │  ┌─────────────────────────────────────────────────────┐   ││
│  │  │            gRPC Stream Handler (per client)         │   ││
│  │  │         StreamPrices (bidirectional streaming)      │   ││
│  │  └─────────────────────────────────────────────────────┘   ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Message Flow

### Price Update Flow
```
Binance WS → PriceUpdate → PriceDistributor → [All Active Streams]
                                          ↓
                              IndicatorEngine (async calc)
                                          ↓
                              IndicatorUpdate → [All Active Streams]
```

### Client Request Flow
```
Client gRPC Stream → IndicatorRequest → StreamRegistry
                                            ↓
                                    IndicatorEngine (config update)
```

## Error Handling Strategy

| Scenario | Handling |
|----------|----------|
| Binance WS disconnect | Auto-reconnect with exponential backoff |
| gRPC client disconnect | Clean up stream from registry |
| Indicator calculation error | Log and skip update, don't crash stream |
| Shutdown timeout | Force close streams after 5s |
| Memory pressure | Configurable price history limits |
