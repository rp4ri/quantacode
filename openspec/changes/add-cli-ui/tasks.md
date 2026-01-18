## 1. CLI Framework Setup
- [x] 1.1 Create cmd/cli/main.go file
- [x] 1.2 Initialize Cobra root command
- [x] 1.3 Add 'chat' subcommand
- [x] 1.4 Add --server flag (default: localhost:50051)
- [x] 1.5 Add gRPC client initialization with server address
- [x] 1.6 Print connection status on startup
- [x] 1.7 Handle connection errors gracefully

## 2. Bubble Tea TUI Framework
- [x] 2.1 Initialize Bubble Tea program
- [x] 2.2 Define TUI model struct with state
- [x] 2.3 Implement Init, Update, View methods
- [x] 2.4 Add message types for price updates, chat messages, typing indicator

## 3. Price Ticker Component
- [x] 3.1 Create price ticker display component
- [x] 3.2 Format price display with symbol, price, change
- [x] 3.3 Implement green/red color coding for price changes
- [x] 3.4 Update ticker on price update messages

## 4. Chat Interface
- [x] 4.1 Create internal/ui/chat/view.go file
- [x] 4.2 Implement split-screen layout (ticker top, chat bottom)
- [x] 4.3 Create chat message display with scrollable history
- [x] 4.4 Add user input textarea at bottom
- [x] 4.5 Implement question submission via Enter key
- [x] 4.6 Send questions to AI (implemented client-side via OpenRouter)
- [x] 4.7 Display AI responses in chat area
- [x] 4.8 Add styling with Lipgloss (user messages right-aligned, AI left-aligned)

## 5. Typing Indicator
- [x] 5.1 Create typing indicator animation (dots: ., .., ...)
- [x] 5.2 Show indicator when waiting for AI response
- [x] 5.3 Hide indicator when AI response received
- [x] 5.4 Handle indicator state in Bubble Tea model

## 6. Indicators Panel
- [x] 6.1 Create internal/ui/indicators/panel.go file
- [x] 6.2 Implement right sidebar layout
- [x] 6.3 Display RSI, SMA, EMA values with labels
- [x] 6.4 Implement color coding for RSI (red >70, blue <30, white else)
- [x] 6.5 Update indicators on price/indicator update messages
- [x] 6.6 Add styling with Lipgloss

## 7. gRPC Client Integration
- [x] 7.1 Create gRPC client stub for MarketDataService
- [x] 7.2 Implement StreamPrices subscription
- [x] 7.3 Handle incoming PriceUpdate messages
- [x] 7.4 Implement AI analysis for user questions (client-side via OpenRouter)
- [x] 7.5 Handle streaming AI responses (SSE streaming implemented)
- [x] 7.6 Map gRPC messages to Bubble Tea messages

## 8. TUI Layout Composition
- [x] 8.1 Compose final layout: top ticker, middle chat, bottom input, right indicators
- [x] 8.2 Implement responsive layout handling
- [x] 8.3 Add keyboard shortcuts (q/ctrl+c to quit)
- [x] 8.4 Handle window resize events

## 9. Validation
- [x] 9.1 Test CLI with `quantacode chat --server=localhost:50051`
- [x] 9.2 Verify price ticker updates in real-time
- [x] 9.3 Test user question submission and AI response display
- [x] 9.4 Verify typing indicator appears/disappears correctly
- [x] 9.5 Test indicator panel color coding (RSI >70, <30)
- [x] 9.6 Test connection error handling
- [x] 9.7 Run `go test ./cmd/cli/...` and `go test ./internal/ui/...`
- [x] 9.8 Test keyboard shortcuts for quit
