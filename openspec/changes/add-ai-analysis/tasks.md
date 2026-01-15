## 1. OpenRouter Client Setup
- [ ] 1.1 Create internal/infra/openrouter/client.go file
- [ ] 1.2 Implement HTTP client with request/response structures
- [ ] 1.3 Add Bearer token authentication from OPENROUTER_API_KEY env var
- [ ] 1.4 Configure deepseek/deepseek-chat model endpoint
- [ ] 1.5 Implement chat completion request/response parsing
- [ ] 1.6 Add unit tests for client (mock HTTP responses)

## 2. Analysis Context Builder
- [ ] 2.1 Create internal/domain/analysis/context.go file
- [ ] 2.2 Define AnalysisContext struct with price, indicators, trends
- [ ] 2.3 Implement format function to create structured prompt
- [ ] 2.4 Include current price, RSI, SMA, EMA values
- [ ] 2.5 Add "explain why" instruction to prompt
- [ ] 2.6 Add unit tests for context formatting

## 3. Protocol Buffer Updates
- [ ] 3.1 Add AIAnalysis message to proto/market_data.proto
- [ ] 3.2 Add StreamAIAnalysis method to gRPC service
- [ ] 3.3 Regenerate Go code with protoc

## 4. gRPC Server Integration
- [ ] 4.1 Implement AI analysis trigger logic (price change >1%)
- [ ] 4.2 Create OpenRouter client instance in gRPC server
- [ ] 4.3 Build analysis context from current indicators
- [ ] 4.4 Stream AI responses to gRPC client
- [ ] 4.5 Add "Not financial advice" disclaimer to responses

## 5. Rate Limiting & Error Handling
- [ ] 5.1 Implement exponential backoff for OpenRouter rate limits (429 status)
- [ ] 5.2 Add retry logic with max attempts
- [ ] 5.3 Handle OpenRouter API errors gracefully
- [ ] 5.4 Add logging for AI analysis requests

## 6. Validation
- [ ] 6.1 Test OpenRouter client with real API key
- [ ] 6.2 Verify analysis context formatting produces valid prompts
- [ ] 6.3 Test gRPC streaming of AI responses
- [ ] 6.4 Verify rate limiting behavior (simulate 429 responses)
- [ ] 6.5 Check that 1% price change trigger works correctly
- [ ] 6.6 Run all unit tests
- [ ] 6.7 Verify financial disclaimers present in responses
