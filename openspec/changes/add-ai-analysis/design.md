## Context
The QuantaCode CLI needs to provide actionable trading insights, not just raw data. AI-powered analysis via DeepSeek can interpret technical indicators, explain trends, and suggest strategies conversationally. This requires integrating OpenRouter API with rate limiting, structured prompt engineering, and streaming responses to gRPC clients.

## Goals / Non-Goals
- Goals: AI-generated trading analysis, conversational insights explaining "why not what", real-time streaming of responses
- Non-Goals: Trading execution, portfolio management, historical analysis, sentiment analysis from social media

## Decisions

### OpenRouter vs Direct DeepSeek API
- **Decision**: Use OpenRouter API gateway
- **Rationale**: Unified API for multiple models, simpler key management, free tier support
- **Alternative**: Direct DeepSeek API - rejected due to complexity of managing multiple providers

### Analysis Trigger: 1% Price Change Threshold
- **Decision**: Trigger AI analysis only when price changes >1%
- **Rationale**: Reduces API usage (1000/day limit), avoids spamming user with minor fluctuations
- **Alternative**: Every N updates - rejected as less meaningful for trading analysis
- **Alternative**: User-configurable threshold - rejected for simplicity (can add later)

### Structured Prompt Engineering
- **Decision**: Build structured prompt with price, indicators, trends, and explicit "explain why" instruction
- **Rationale**: Ensures AI provides actionable insights rather than generic descriptions
- **Alternative**: Unstructured prompt - rejected for inconsistent output quality

### Streaming vs Batch Responses
- **Decision**: Stream AI responses via gRPC to CLI client
- **Rationale**: Real-time user experience, progressive rendering in TUI
- **Alternative**: Wait for complete response - rejected as poor UX for long analyses

### Rate Limiting Strategy
- **Decision**: Exponential backoff with max 3 retries on HTTP 429
- **Rationale**: Standard pattern for API rate limits, graceful degradation
- **Alternative**: Queue requests - rejected due to complexity and latency

### Financial Disclaimers
- **Decision**: Append "Not financial advice" to all AI responses
- **Rationale**: Regulatory requirement, risk management, user protection
- **Implementation**: Add as footer to gRPC StreamAIAnalysis responses

## Data Structures

```go
// OpenRouter Client
type OpenRouterClient struct {
    apiKey    string
    httpClient *http.Client
    baseURL   string
}

type ChatRequest struct {
    Model    string   `json:"model"`
    Messages []Message `json:"messages"`
}

type ChatResponse struct {
    Choices []Choice `json:"choices"`
}

// Analysis Context
type AnalysisContext struct {
    Symbol     string
    Price      float64
    PriceChange float64
    Indicators IndicatorValues
    Trend      string
    Timestamp  time.Time
}

type IndicatorValues struct {
    RSI   float64
    SMA14 float64
    EMA12 float64
}

// Proto Messages
message AIAnalysis {
    string symbol = 1;
    string analysis = 2;  // AI-generated text
    int64 timestamp = 3;
    bool is_streaming = 4;  // Last message has this false
}

service MarketDataService {
    rpc StreamAIAnalysis(PriceUpdate) returns (stream AIAnalysis);
}
```

## Rate Limiting Algorithm

```go
// Exponential backoff: 1s, 2s, 4s, 8s (max 4 attempts)
func (c *OpenRouterClient) retryWithBackoff(fn func() error) error {
    maxRetries := 3
    for i := 0; i <= maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        if !isRateLimitError(err) {
            return err
        }
        if i == maxRetries {
            return err
        }
        delay := time.Duration(1<<uint(i)) * time.Second
        time.Sleep(delay)
    }
    return nil
}
```

## Risks / Trade-offs

### OpenRouter API Limits
- **Risk**: 1000 requests/day free tier may be insufficient for active trading
- **Mitigation**: 1% price change threshold reduces calls significantly; recommend upgrade for power users
- **Trade-off**: Conservative usage limit vs. more frequent analysis

### AI Response Latency
- **Risk**: DeepSeek API response time 1-3s affects real-time UX
- **Mitigation**: Stream responses progressively, show "analyzing..." indicator
- **Trade-off**: Accept 1-3s latency for high-quality insights

### Cost Management
- **Risk**: Users might exceed free tier quickly
- **Mitigation**: Clear documentation of limits, count API calls in logs
- **Trade-off**: Free tier limitation vs. paid tier recommendation

### Prompt Injection
- **Risk**: Malicious input from price data could manipulate AI responses
- **Mitigation**: Sanitize numeric data, validate indicator ranges, use system prompt to enforce constraints
- **Trade-off**: Minimal risk since input is structured numeric data

### Hallucinations
- **Risk**: AI may provide incorrect or hallucinated financial advice
- **Mitigation**: Strong system prompt emphasizing "not financial advice", instruct to focus on data interpretation
- **Trade-off**: Accept minor hallucination risk with disclaimers vs. no AI at all

## Error Handling Strategy

### OpenRouter API Failures
- 429 Too Many Requests: Retry with exponential backoff
- 401 Unauthorized: Log error, return authentication error to client
- 500 Server Error: Retry once, then return error to client
- Network Timeout: Retry with backoff

### gRPC Streaming Errors
- Client disconnect: Close stream, cleanup resources
- OpenRouter failure during stream: Send error message to client, close stream gracefully

## Migration Plan
N/A - new capability, no existing code to migrate

## Open Questions
- Should we cache AI responses for identical market conditions?
- Should users be able to customize the analysis prompt?
- Should we support additional models (e.g., GPT-4, Claude) via OpenRouter?
- How should we handle cases where price never changes >1% for extended periods?
