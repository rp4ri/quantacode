# Tasks: gRPC Streaming Server Implementation

## Phase 1: Protocol Buffer Definitions

### Task 1.1: Define proto schema extensions
- [x] **Description**: Add `IndicatorRequest` and extend messages in `proto/marketdata.proto`
- **Files**: proto/marketdata.proto (new or existing)
- **Validation**: `protoc --go_out=. --go-grpc_out=. proto/marketdata.proto` compiles without errors
- **Dependencies**: None
- **Estimated**: 30min

### Task 1.2: Generate Go code from proto
- [x] **Description**: Run protoc to generate Go structs and gRPC stubs
- **Files**: proto/marketdata.pb.go, proto/marketdata_grpc.pb.go
- **Validation**: Generated files compile, no missing dependencies
- **Dependencies**: Task 1.1
- **Estimated**: 15min

## Phase 2: Core Server Infrastructure

### Task 2.1: Create cmd/server/main.go
- [x] **Description**: Implement server entrypoint with gRPC server initialization on :50051
- **Files**: cmd/server/main.go
- **Validation**: `go build -o bin/server ./cmd/server` succeeds, server starts on :50051
- **Dependencies**: Phase 1
- **Estimated**: 45min

### Task 2.2: Implement signal handler
- [x] **Description**: Add SIGTERM/SIGINT handling with context cancellation
- **Files**: cmd/server/main.go
- **Validation**: Server responds to kill signals, logs shutdown initiation
- **Dependencies**: Task 2.1
- **Estimated**: 20min

### Task 2.3: Create gRPC server handler struct
- [x] **Description**: Implement `internal/grpc/server/handler.go` with StreamPrices method signature
- **Files**: internal/grpc/server/handler.go
- **Validation**: `go build` succeeds, handler struct compiles
- **Dependencies**: Task 2.1
- **Estimated**: 30min

## Phase 3: Binance Integration

### Task 3.1: Create Binance WebSocket client wrapper
- [x] **Description**: Implement connection management in `internal/infra/binance/client.go`
- **Files**: internal/infra/binance/client.go
- **Validation**: Connects to Binance, receives ticker messages, parses JSON correctly
- **Dependencies**: Task 2.1
- **Estimated**: 1h

### Task 3.2: Implement reconnection logic
- [x] **Description**: Add exponential backoff reconnection with max retries
- **Files**: internal/infra/binance/client.go
- **Validation**: Simulate disconnection, verify reconnection within 30s max
- **Dependencies**: Task 3.1
- **Estimated**: 30min

### Task 3.3: Create price distributor
- [x] **Description**: Implement fan-out mechanism to distribute price updates to multiple streams
- **Files**: internal/grpc/server/handler.go
- **Validation**: Multiple gRPC clients receive same price updates
- **Dependencies**: Task 2.3, Task 3.1
- **Estimated**: 45min

## Phase 4: Indicator Calculations

### Task 4.1: Implement RSI calculator
- [x] **Description**: Add RSI calculation in `internal/domain/indicators/rsi.go`
- **Files**: internal/domain/indicators/rsi.go
- **Validation**: RSI values match expected results for known input data
- **Dependencies**: None
- **Estimated**: 45min

### Task 4.2: Implement SMA calculator
- [x] **Description**: Add SMA calculation in `internal/domain/indicators/sma.go`
- **Files**: internal/domain/indicators/sma.go
- **Validation**: SMA values match expected results for known input data
- **Dependencies**: None
- **Estimated**: 30min

### Task 4.3: Implement EMA calculator
- [x] **Description**: Add EMA calculation in `internal/domain/indicators/ema.go`
- **Files**: internal/domain/indicators/ema.go
- **Validation**: EMA values match expected results for known input data
- **Dependencies**: None
- **Estimated**: 30min

### Task 4.4: Create indicator engine
- [x] **Description**: Implement centralized indicator calculation and subscription management
- **Files**: internal/grpc/server/indicator_engine.go
- **Validation**: Multiple indicators calculate correctly for same price stream
- **Dependencies**: Task 4.1, Task 4.2, Task 4.3
- **Estimated**: 1h

## Phase 5: Bidirectional Streaming

### Task 5.1: Implement StreamPrices server streaming
- [x] **Description**: Complete `StreamPrices` method for server-side price updates
- **Files**: internal/grpc/server/handler.go
- **Validation**: Clients connect and receive continuous price updates
- **Dependencies**: Task 2.3, Task 3.3
- **Estimated**: 45min

### Task 5.2: Implement IndicatorRequest handling
- [x] **Description**: Process client requests to subscribe/unsubscribe from indicators
- **Files**: internal/grpc/server/handler.go
- **Validation**: Client can dynamically add/remove indicator subscriptions
- **Dependencies**: Task 5.1
- **Estimated**: 45min

### Task 5.3: Interleave price and indicator updates
- [x] **Description**: Send combined stream of PriceUpdate and IndicatorUpdate messages
- **Files**: internal/grpc/server/handler.go
- **Validation**: Clients receive both price and indicator updates on same stream
- **Dependencies**: Task 5.2, Task 4.4
- **Estimated**: 30min

## Phase 6: Graceful Shutdown

### Task 6.1: Implement stream draining
- [x] **Description**: Track active streams, notify clients, drain within timeout
- **Files**: internal/grpc/server/handler.go, cmd/server/main.go
- **Validation**: Active streams receive notification, complete within 5s
- **Dependencies**: Task 2.2, Task 5.1
- **Estimated**: 45min

### Task 6.2: Implement Binance connection close
- [x] **Description**: Cleanly close WebSocket during shutdown
- **Files**: internal/infra/binance/client.go
- **Validation**: Connection closes gracefully, no orphaned connections
- **Dependencies**: Task 3.1, Task 6.1
- **Estimated**: 20min

### Task 6.3: Implement gRPC server graceful stop
- [x] **Description**: Use GracefulStop() with 5s timeout
- **Files**: cmd/server/main.go
- **Validation**: Server stops accepting connections, in-flight requests complete
- **Dependencies**: Task 6.1
- **Estimated**: 30min

### Task 6.4: Integrate shutdown sequence
- [x] **Description**: Coordinate signal handler, stream drain, Binance close, server stop
- **Files**: cmd/server/main.go
- **Validation**: Full shutdown sequence completes within 10s, clean exit
- **Dependencies**: Task 6.1, Task 6.2, Task 6.3
- **Estimated**: 30min

## Phase 7: Testing and Validation

### Task 7.1: Write unit tests for indicators
- [x] **Description**: Add table-driven tests for RSI, SMA, EMA calculations
- **Files**: internal/domain/indicators/*_test.go
- **Validation**: `go test ./internal/domain/indicators/... -v` passes
- **Dependencies**: Phase 4
- **Estimated**: 1h

### Task 7.2: Write integration test for gRPC streaming
- [x] **Description**: Create test that starts server, connects client, receives updates (manual testing completed)
- **Files**: internal/grpc/server/handler_test.go
- **Validation**: `go test ./internal/grpc/server/... -v` passes
- **Dependencies**: Phase 5
- **Estimated**: 1h

### Task 7.3: Test graceful shutdown
- [x] **Description**: Verify shutdown sequence with active connections
- **Files**: cmd/server/main_test.go (or manual test script)
- **Validation**: Server shuts down cleanly with active clients, no data loss
- **Dependencies**: Phase 6
- **Estimated**: 30min

### Task 7.4: Run linter and fix issues
- [x] **Description**: Run golangci-lint and resolve all warnings
- **Files**: All Go files
- **Validation**: `golangci-lint run ./...` passes with no errors
- **Dependencies**: All previous tasks
- **Estimated**: 30min

## Summary

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1: Protocol Buffers | 2 | 45min |
| Phase 2: Server Infrastructure | 3 | 1h 35min |
| Phase 3: Binance Integration | 3 | 2h 15min |
| Phase 4: Indicator Calculations | 4 | 2h 45min |
| Phase 5: Bidirectional Streaming | 3 | 2h |
| Phase 6: Graceful Shutdown | 4 | 2h 5min |
| Phase 7: Testing and Validation | 4 | 3h 30min |
| **Total** | **23** | **~15 hours** |

## Parallelization Opportunities

- **Phase 1** can run in parallel with any work
- **Tasks 4.1, 4.2, 4.3** (individual indicators) can run in parallel
- **Phase 7** (testing) can begin after each phase completes
- **Phase 2.3** (handler struct) can start before **Phase 1** completes (interface-based)
