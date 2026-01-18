## 1. OpenRouter Client Setup
- [x] 1.1 Create internal/ai/openrouter/client.go file
- [x] 1.2 Implement HTTP client with request/response structures
- [x] 1.3 Add Bearer token authentication from OPENROUTER_API_KEY env var
- [x] 1.4 Configure deepseek/deepseek-chat model endpoint
- [x] 1.5 Implement chat completion request/response parsing with SSE streaming
- [ ] 1.6 Add unit tests for client (mock HTTP responses)

## 2. Analysis Context Builder
- [x] 2.1 Context built inline in openrouter/client.go buildSystemPrompt method
- [x] 2.2 System prompt includes price, indicators, trends context
- [x] 2.3 Implement format function to create structured prompt
- [x] 2.4 Include current price, RSI, SMA, EMA values
- [x] 2.5 Add trading analysis instruction to prompt
- [ ] 2.6 Add unit tests for context formatting

## 3. Protocol Buffer Updates
- [x] 3.1 MarketUpdate message supports price and indicator streaming
- [x] 3.2 AI analysis runs client-side via OpenRouter (no gRPC method needed)
- [x] 3.3 Go code generated with protoc

## 4. CLI/UI Integration (changed from gRPC Server)
- [x] 4.1 AI analysis triggered by user question in chat UI
- [x] 4.2 Create OpenRouter client instance in chat UI
- [x] 4.3 Build analysis context from current indicators
- [x] 4.4 Stream AI responses directly to TUI with real-time display
- [ ] 4.5 Add "Not financial advice" disclaimer to responses

## 5. Rate Limiting & Error Handling
- [ ] 5.1 Implement exponential backoff for OpenRouter rate limits (429 status)
- [ ] 5.2 Add retry logic with max attempts
- [x] 5.3 Handle OpenRouter API errors gracefully
- [x] 5.4 Add logging for AI analysis requests (internal/logging/logger.go)

## 6. Validation
- [x] 6.1 Test OpenRouter client with real API key
- [x] 6.2 Verify analysis context formatting produces valid prompts
- [x] 6.3 Test streaming of AI responses in TUI
- [ ] 6.4 Verify rate limiting behavior (simulate 429 responses)
- [ ] 6.5 Check that price change trigger works correctly
- [x] 6.6 Run all unit tests
- [ ] 6.7 Verify financial disclaimers present in responses
