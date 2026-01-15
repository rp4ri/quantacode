# Change: Add AI-Powered Market Analysis

## Why
Integrate DeepSeek AI via OpenRouter to provide conversational trading analysis that explains market conditions, interprets technical indicators, and suggests strategies based on real-time price data, moving beyond raw data to actionable insights.

## What Changes
- Create OpenRouter HTTP client for https://openrouter.ai/api/v1/chat/completions
- Configure deepseek/deepseek-chat model with Bearer authentication from environment variable
- Implement analysis context builder that formats price + indicators into structured prompt for AI
- Include current values, trends, and "explain why" instruction in prompts
- Add AIAnalysis message to proto definitions
- Server calls OpenRouter when price changes >1%, streams AI response via gRPC
- Handle OpenRouter rate limiting with exponential backoff
- Add financial disclaimers to AI responses

## Impact
- Affected specs: openrouter-client, analysis-context, grpc-ai-streaming
- New files: internal/infra/openrouter/client.go, internal/domain/analysis/context.go
- Modified files: proto/market_data.proto, internal/grpc server
- New capability: AI-driven market analysis with conversational insights
