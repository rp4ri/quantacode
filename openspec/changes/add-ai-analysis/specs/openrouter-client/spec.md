## ADDED Requirements

### Requirement: OpenRouter HTTP Client
The system SHALL provide an HTTP client for the OpenRouter API endpoint https://openrouter.ai/api/v1/chat/completions.

#### Scenario: Initialize client with API key
- **WHEN** OpenRouterClient is created
- **THEN** it loads OPENROUTER_API_KEY from environment variable

#### Scenario: Send chat completion request
- **WHEN** client sends request to chat completions endpoint
- **THEN** it uses Bearer token authentication

### Requirement: DeepSeek Model Configuration
The system SHALL use the deepseek/deepseek-chat model for trading analysis requests.

#### Scenario: Configure model parameter
- **WHEN** chat request is constructed
- **THEN** model field is set to "deepseek/deepseek-chat"

### Requirement: Bearer Authentication
The system SHALL authenticate OpenRouter API requests using Bearer token from OPENROUTER_API_KEY environment variable.

#### Scenario: Add authorization header
- **WHEN** HTTP request is sent to OpenRouter
- **THEN** Authorization header contains "Bearer <API_KEY>"

### Requirement: Rate Limiting with Exponential Backoff
The system SHALL implement exponential backoff retry logic when OpenRouter returns HTTP 429 (rate limit exceeded).

#### Scenario: Retry on rate limit
- **WHEN** OpenRouter returns 429 status code
- **THEN** client retries with 1s, 2s, 4s delays (max 3 retries)

#### Scenario: Max retries exhausted
- **WHEN** all 3 retry attempts fail with 429
- **THEN** client returns error to caller

### Requirement: Error Handling
The system SHALL handle OpenRouter API errors gracefully without crashing the application.

#### Scenario: Authentication error
- **WHEN** OpenRouter returns 401 Unauthorized
- **THEN** error is logged and returned to caller with clear message

#### Scenario: Server error
- **WHEN** OpenRouter returns 5xx status code
- **THEN** error is logged and retry is attempted once
