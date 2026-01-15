# OpenSpec Proposals for QuantaCode

## Phase 1: Foundation

1. "Initialize Go module with go.mod, create basic project structure with cmd/server, cmd/cli, internal/domain, internal/infra folders. Add .gitignore for Go. Setup README with project description."

2. "Define proto/market_data.proto with PriceUpdate message (symbol, price, volume, timestamp). Generate Go code with protoc. Create basic gRPC service definition for StreamPrices."

3. "Implement internal/infra/binance/client.go: connect to wss://stream.binance.com:9443/ws/btcusdt@ticker, parse JSON ticker, return PriceUpdate struct. Handle reconnection on error."

## Phase 2: Indicators

4. "Create internal/domain/indicators/rsi.go: implement RSI calculation with period parameter. Use circular buffer for price history. Add unit tests with known input/output values."

5. "Add internal/domain/indicators/sma.go and ema.go: implement Simple and Exponential Moving Average. Both accept period parameter. Include table-driven unit tests."

6. "Create file to aggregate all indicators, maintains price buffer, updates all indicators on new price. Return struct with all values."

## Phase 3: gRPC Server

7. "Implement cmd/server/main.go: start gRPC server on :50051. Create internal/grpc/server/handler.go with StreamPrices method that subscribes to Binance and streams PriceUpdate."

8. "Add bidirectional streaming: extend proto with IndicatorRequest message. Server streams both PriceUpdate and calculated indicators (RSI, SMA, EMA) to clients."

9. "Add graceful shutdown: handle SIGTERM/SIGINT in server, close Binance connection, drain active streams, shutdown gRPC server with 5s timeout."

## Phase 4: AI Integration

10. "Create internal/infra/openrouter/client.go: HTTP client for https://openrouter.ai/api/v1/chat/completions. Use deepseek/deepseek-chat model. Handle Bearer auth from env var."

11. "Implement internal/domain/analysis/context.go: format price + indicators into structured prompt for AI. Include current values, trends, and 'explain why' instruction."

12. "Add AIAnalysis message to proto. Server calls OpenRouter when price changes >1%, streams AI response. Handle rate limiting with exponential backoff."

## Phase 5: CLI Client

13. "Create cmd/cli/main.go with Cobra: root command, 'chat' subcommand. Connect to gRPC server via flag --server=localhost:50051. Print connection status."

14. "Implement internal/ui/chat/view.go with Bubble Tea: split screen with price ticker on top, chat messages below. Use lipgloss for styling (green/red for price changes)."

15. "Add user input: textarea at bottom for questions. Send to gRPC server, display AI response in chat. Show typing indicator while AI generates response."

16. "Add internal/ui/indicators/panel.go: right sidebar showing RSI, SMA, EMA values. Color-code RSI (red >70, blue <30, white else). Update in real-time."

## Phase 6: Polish

17. "Add cmd/cli flags: --symbol (default BTCUSDT), --openrouter-key, --help. Load from env vars if not provided. Validate required flags at startup."

18. "Create Dockerfile for server: multi-stage build, copy binary, expose port 50051. Add docker-compose.yml for local development with server + sample CLI."

19. "Setup .github/workflows/test.yml: run tests and golangci-lint on push to dev. Add workflow to build and push Docker image on push to prod."

