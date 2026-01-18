# Project Context

## Purpose
QuantaCode is an AI-powered trading analysis CLI that provides real-time market insights through conversational interaction. The tool streams live cryptocurrency data from Binance, calculates technical indicators in real-time, and leverages DeepSeek 4.1 via OpenRouter to generate actionable trading analysis. Unlike traditional CLI tools, it offers a Claude Code-style conversational interface where users chat with an AI analyst that interprets market conditions, explains indicator signals, and suggests trading strategies based on live data.

## Tech Stack
- **Language**: Go 1.21+
- **CLI Framework**: Cobra (command parsing) + Bubble Tea (TUI framework)
- **Styling**: Lipgloss (terminal styling)
- **RPC Protocol**: gRPC with Protocol Buffers
- **AI Integration**: OpenRouter API (DeepSeek 4.1)
- **Market Data**: Binance WebSocket API (real-time streams)
- **Infrastructure**: OpenTofu (Terraform alternative)
- **CI/CD**: GitHub Actions
- **Deployment**: Railway (gRPC server hosting)

## Project Structure

```
quantacode/
├── .github/
│   └── workflows/          # GitHub Actions CI/CD
├── .opencode/              # OpenCode AI context
├── .windsurf/              # Windsurf AI context
├── .claude/                # Claude AI context
├── .agent/                 # Agent configuration
├── openspec/               # OpenSpec files
├── cmd/
│   ├── server/             # gRPC server entrypoint
│   │   └── main.go
│   └── cli/                # CLI client entrypoint
│       └── main.go
├── internal/
│   ├── domain/             # Business logic
│   │   ├── indicators/     # Technical indicator calculations
│   │   └── analysis/       # Market analysis models
│   ├── infra/              # External integrations
│   │   ├── binance/        # Binance API client
│   │   └── openrouter/     # OpenRouter AI client
│   ├── grpc/               # gRPC implementation
│   │   ├── server/         # Server handlers
│   │   └── client/         # Client stubs
│   └── ui/                 # Bubble Tea components
│       ├── chat/           # Chat interface
│       └── indicators/     # Indicator display
├── pkg/                    # Shared utilities (optional)
├── terraform/              # OpenTofu infrastructure
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Project Conventions

### Code Style
- Follow Go's official style guide and `gofmt` formatting
- Use meaningful variable names: `priceUpdate` not `pu`
- Package naming: lowercase, single-word (e.g., `indicators`, `grpc`)
- Error handling: Always check errors, wrap with `fmt.Errorf("operation failed: %w", err)`
- Comments: Godoc format for all exported functions/types

### Architecture Patterns
- **Clean Architecture**: Domain logic separated from infrastructure
  - `internal/domain`: Business logic (indicators, analysis)
  - `internal/infra`: External services (Binance, OpenRouter)
  - `internal/grpc`: gRPC server/client
  - `internal/ui`: Bubble Tea TUI components
- **Event-driven**: gRPC bidirectional streaming
- **Dependency Injection**: Pass dependencies via constructors
- **Context Propagation**: Use `context.Context` for cancellation

### Testing Strategy
- Place all test files under a dedicated `tests/` folder within each package/module (e.g., `internal/domain/indicators/tests`), using the external package convention where practical.
- Unit tests for indicator calculations (table-driven tests)
- Integration tests for gRPC endpoints
- Mock external APIs using interfaces
- Run tests in CI on every PR to `dev`

### Git Workflow
- **Branching**: `dev` (default), `prod`
- **Commit Convention**: Conventional Commits
  - `feat(grpc): add streaming price updates`
  - `fix(indicators): correct RSI calculation`
- **PR Requirements**: 
  - Passing tests + linter (`golangci-lint`)
  - Merge to `dev` for development
  - Merge `dev` → `prod` for releases
- **Releases**: Semantic versioning (v1.2.3), automated via GitHub Actions

## Domain Context

### Trading Indicators
Real-time calculations:
- **RSI**: Momentum oscillator (overbought >70, oversold <30)
- **SMA/EMA**: Trend indicators
- **MACD**: Momentum + trend
- **Bollinger Bands**: Volatility indicator

### Market Data
- **Price Updates**: 100ms intervals for active pairs
- **Order Book**: Top 20 bid/ask levels
- **Timeframes**: 1m, 5m, 15m, 1h, 4h, 1d

### AI Analysis
DeepSeek receives:
- Current price action (trend, momentum, volatility)
- Technical indicator values
- Volume and order book data
AI responds conversationally, explaining "why" not just "what"

## Important Constraints

### Technical
- **Binance Rate Limits**: 1200 req/min (REST), 5 connections (WebSocket)
- **gRPC Message Size**: Max 4MB per message
- **Terminal Compatibility**: Linux, macOS, Windows (WSL)
- **Memory**: <100MB CLI client, <500MB server

### Business
- **No Trading**: Analysis only, no order execution
- **Latency SLA**: Updates render within 200ms
- **Cost**: OpenRouter calls capped at 1000/day (free tier)

### Regulatory
- **Disclaimers**: "Not financial advice" on all AI output
- **Privacy**: No user data stored (stateless)
- **API Keys**: User-provided OpenRouter keys

## External Dependencies

### Binance API
- **WebSocket**: `wss://stream.binance.com:9443/stream`
  - `@ticker`, `@depth`, `@kline` streams
- **REST**: `https://api.binance.us/api/v3`
- **Auth**: None required for public data
- **Docs**: https://binance-docs.github.io/apidocs/spot/en/

### OpenRouter API
- **Endpoint**: `https://openrouter.ai/api/v1/chat/completions`
- **Model**: DeepSeek 4.1 (`deepseek/deepseek-chat`)
- **Auth**: Bearer token via `Authorization` header
- **Docs**: https://openrouter.ai/docs

### Infrastructure
- **Railway**: gRPC server hosting
- **GitHub Actions**: CI/CD for `dev` → `prod` pipeline
- **OpenTofu**: Infrastructure as Code
- **GHCR**: Container registry for server images

### Go Dependencies
- `github.com/spf13/cobra`: CLI framework
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Terminal styling
- `google.golang.org/grpc`: gRPC
- `google.golang.org/protobuf`: Protocol Buffers
- `github.com/gorilla/websocket`: Binance WebSocket client
- `github.com/stretchr/testify`: Testing
