## 1. Project Initialization
- [x] 1.1 Initialize Go module with `go mod init github.com/rp4ri/quantacode`
- [ ] 1.2 Create directory structure: cmd/server, cmd/cli, internal/domain, internal/infra
- [ ] 1.3 Create Go .gitignore (bin/, pkg/, .env, vendor/)
- [ ] 1.4 Create README.md with project description

## 2. Protocol Buffers Setup
- [ ] 2.1 Create proto/market_data.proto with PriceUpdate message
- [ ] 2.2 Define StreamPrices gRPC service in proto
- [ ] 2.3 Generate Go code with protoc and protoc-gen-go
- [ ] 2.4 Add protobuf dependencies to go.mod

## 3. Binance WebSocket Client
- [ ] 3.1 Create internal/infra/binance/client.go file
- [ ] 3.2 Implement WebSocket connection to wss://stream.binance.com:9443/ws/btcusdt@ticker
- [ ] 3.3 Define PriceUpdate struct for parsed data
- [ ] 3.4 Implement JSON ticker parsing from Binance WebSocket
- [ ] 3.5 Add reconnection logic with exponential backoff
- [ ] 3.6 Add gorilla/websocket dependency to go.mod

## 4. gRPC Service Integration
- [ ] 4.1 Create internal/grpc/server/market_data_server.go
- [ ] 4.2 Implement StreamPrices gRPC method
- [ ] 4.3 Integrate Binance client with gRPC streaming
- [ ] 4.4 Add gRPC server dependencies to go.mod

## 5. Validation
- [ ] 5.1 Test WebSocket connection and message parsing
- [ ] 5.2 Verify gRPC code generation
- [ ] 5.3 Run `go mod tidy` to ensure dependencies
- [ ] 5.4 Compile project to verify no syntax errors
