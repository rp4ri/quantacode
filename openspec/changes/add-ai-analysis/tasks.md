## 1. OpenRouter Client Setup
- [x] 1.1 Create internal/ai/openrouter/client.go file
- [x] 1.2 Implement HTTP client with request/response structures
- [x] 1.3 Add Bearer token authentication from OPENROUTER_API_KEY env var
- [x] 1.4 Configure deepseek/deepseek-chat model endpoint
- [x] 1.5 Implement chat completion request/response parsing with SSE streaming
- [x] 1.6 Add IndicatorHistory struct for passing historical data to AI

## 2. Analysis Context Builder
- [x] 2.1 Context built inline in openrouter/client.go buildSystemPrompt method
- [x] 2.2 System prompt includes price, indicators, trends context
- [x] 2.3 Implement format function to create structured prompt
- [x] 2.4 Include current price, RSI, SMA, EMA values
- [x] 2.5 Add trading analysis instruction to prompt
- [x] 2.6 Include full indicator history (last 30 candles) in prompt as markdown table

## 3. Protocol Buffer Updates
- [x] 3.1 MarketUpdate message supports price and indicator streaming
- [x] 3.2 AI analysis runs client-side via OpenRouter (no gRPC method needed)
- [x] 3.3 Go code generated with protoc
- [x] 3.4 Added indicator history fields to IndicatorUpdate message

## 4. CLI/UI Integration
- [x] 4.1 AI analysis triggered by user question in chat UI
- [x] 4.2 Create OpenRouter client instance in chat UI
- [x] 4.3 Build analysis context from current indicators and history
- [x] 4.4 Stream AI responses directly to TUI with real-time display
- [x] 4.5 AI only analyzes when explicitly asked (not on every message)
- [x] 4.6 Added /clear command to clear chat history
- [x] 4.7 Added /pairs command to switch trading pairs
- [x] 4.8 Arrow keys recall previous messages (input history)

## 5. Error Handling & Logging
- [x] 5.1 Handle OpenRouter API errors gracefully
- [x] 5.2 Add logging for AI analysis requests (internal/logging/logger.go)
- [x] 5.3 Logs write to file only (not stdout) to avoid polluting TUI

## 6. Validation
- [x] 6.1 Test OpenRouter client with real API key
- [x] 6.2 Verify analysis context formatting produces valid prompts
- [x] 6.3 Test streaming of AI responses in TUI
- [x] 6.4 Verify indicator history is included in prompts
- [x] 6.5 Test slash commands (/clear, /pairs)
- [x] 6.6 Test input history with arrow keys
