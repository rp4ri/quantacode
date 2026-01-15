## 1. CLI Framework Setup
- [ ] 1.1 Create cmd/cli/main.go file
- [ ] 1.2 Initialize Cobra root command
- [ ] 1.3 Add 'chat' subcommand
- [ ] 1.4 Add --server flag (default: localhost:50051)
- [ ] 1.5 Add gRPC client initialization with server address
- [ ] 1.6 Print connection status on startup
- [ ] 1.7 Handle connection errors gracefully

## 2. Bubble Tea TUI Framework
- [ ] 2.1 Initialize Bubble Tea program
- [ ] 2.2 Define TUI model struct with state
- [ ] 2.3 Implement Init, Update, View methods
- [ ] 2.4 Add message types for price updates, chat messages, typing indicator

## 3. Price Ticker Component
- [ ] 3.1 Create price ticker display component
- [ ] 3.2 Format price display with symbol, price, change
- [ ] 3.3 Implement green/red color coding for price changes
- [ ] 3.4 Update ticker on price update messages

## 4. Chat Interface
- [ ] 4.1 Create internal/ui/chat/view.go file
- [ ] 4.2 Implement split-screen layout (ticker top, chat bottom)
- [ ] 4.3 Create chat message display with scrollable history
- [ ] 4.4 Add user input textarea at bottom
- [ ] 4.5 Implement question submission via Enter key
- [ ] 4.6 Send questions to gRPC server StreamAIAnalysis method
- [ ] 4.7 Display AI responses in chat area
- [ ] 4.8 Add styling with Lipgloss (user messages right-aligned, AI left-aligned)

## 5. Typing Indicator
- [ ] 5.1 Create typing indicator animation (dots: ., .., ...)
- [ ] 5.2 Show indicator when waiting for AI response
- [ ] 5.3 Hide indicator when AI response received
- [ ] 5.4 Handle indicator state in Bubble Tea model

## 6. Indicators Panel
- [ ] 6.1 Create internal/ui/indicators/panel.go file
- [ ] 6.2 Implement right sidebar layout
- [ ] 6.3 Display RSI, SMA, EMA values with labels
- [ ] 6.4 Implement color coding for RSI (red >70, blue <30, white else)
- [ ] 6.5 Update indicators on price/indicator update messages
- [ ] 6.6 Add styling with Lipgloss

## 7. gRPC Client Integration
- [ ] 7.1 Create gRPC client stub for MarketDataService
- [ ] 7.2 Implement StreamPrices subscription
- [ ] 7.3 Handle incoming PriceUpdate messages
- [ ] 7.4 Implement StreamAIAnalysis for user questions
- [ ] 7.5 Handle streaming AIAnalysis responses
- [ ] 7.6 Map gRPC messages to Bubble Tea messages

## 8. TUI Layout Composition
- [ ] 8.1 Compose final layout: top ticker, middle chat, bottom input, right indicators
- [ ] 8.2 Implement responsive layout handling
- [ ] 8.3 Add keyboard shortcuts (q/ctrl+c to quit)
- [ ] 8.4 Handle window resize events

## 9. Validation
- [ ] 9.1 Test CLI with `quantacode chat --server=localhost:50051`
- [ ] 9.2 Verify price ticker updates in real-time
- [ ] 9.3 Test user question submission and AI response display
- [ ] 9.4 Verify typing indicator appears/disappears correctly
- [ ] 9.5 Test indicator panel color coding (RSI >70, <30)
- [ ] 9.6 Test connection error handling
- [ ] 9.7 Run `go test ./cmd/cli/...` and `go test ./internal/ui/...`
- [ ] 9.8 Test keyboard shortcuts for quit
