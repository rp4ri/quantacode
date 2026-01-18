# QuantaCode

QuantaCode is an AI-assisted cryptocurrency trading analysis tool that streams live market data from Binance, computes technical indicators (RSI, SMA, EMA), and provides an interactive terminal UI with AI-powered analysis via OpenRouter.

## Features

- **Real-time price streaming** from Binance (US and global endpoints)
- **Technical indicators**: RSI (14), SMA (14), EMA (14) with 30-candle history
- **AI-powered analysis** using DeepSeek via OpenRouter API
- **Interactive TUI** built with Bubble Tea and Lipgloss
- **15 trading pairs** supported (BTC, ETH, BNB, XRP, ADA, DOGE, SOL, DOT, MATIC, LTC, AVAX, LINK, ATOM, UNI, XLM)
- **Slash commands**: `/clear` to clear chat, `/pairs` to switch trading pairs
- **Input history**: Use arrow keys to recall previous messages

## Architecture

```
┌─────────────────┐     gRPC      ┌─────────────────┐
│   CLI Client    │◄────────────►│   gRPC Server   │
│  (Bubble Tea)   │               │                 │
└────────┬────────┘               └────────┬────────┘
         │                                 │
         │ SSE                             │ WebSocket
         ▼                                 ▼
┌─────────────────┐               ┌─────────────────┐
│   OpenRouter    │               │     Binance     │
│   (AI Analysis) │               │   (Price Data)  │
└─────────────────┘               └─────────────────┘
```

## Prerequisites

- Go 1.21+
- `protoc` (Protocol Buffers compiler)
- OpenRouter API key
- Internet access to Binance (US or global)

### Install protoc (Ubuntu)

```bash
sudo apt update && sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Installation

```bash
git clone <repository>
cd trade-code
go mod tidy
```

## Quick Start

### 1. Start the gRPC Server

```bash
go run ./cmd/server
```

The server will:
- Try to connect to `binance.us` first (for US-based servers)
- Fall back to `binance.com` if US endpoint fails
- Listen on port `50051` by default

Environment variables:
- `PORT`: Server port (default: `50051`)
- `SYMBOL`: Trading symbol (default: `btcusdt`)

### 2. Start the CLI Client

```bash
export OPENROUTER_API_KEY="your-api-key"
go run ./cmd/cli chat
```

Or with flags:

```bash
go run ./cmd/cli chat --server=localhost:50051 --symbol=BTCUSDT --openrouter-key=your-key
```

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--server` | `localhost:50051` | gRPC server address |
| `--symbol` | `BTCUSDT` | Trading pair to subscribe |
| `--openrouter-key` | `$OPENROUTER_API_KEY` | OpenRouter API key |

## Usage

### Keyboard Controls

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `↑` / `↓` | Navigate input history |
| `PgUp` / `PgDown` | Scroll chat history |
| `Esc` / `Ctrl+C` | Exit |

### Slash Commands

| Command | Description |
|---------|-------------|
| `/clear` | Clear chat history |
| `/pairs` | Open trading pair selector |

### AI Analysis

The AI will only provide trading analysis when explicitly asked. Examples:
- "Analiza el mercado"
- "¿Qué señales ves?"
- "Dame tu opinión sobre BTC"

For general questions, it will respond normally without forcing analysis.

## Project Structure

```
trade-code/
├── cmd/
│   ├── cli/          # CLI entrypoint
│   └── server/       # gRPC server entrypoint
├── internal/
│   ├── ai/openrouter/    # OpenRouter client for AI
│   ├── domain/indicators/ # RSI, SMA, EMA calculations
│   ├── grpc/             # gRPC client and server
│   ├── infra/binance/    # Binance WebSocket client
│   ├── logging/          # JSON file logger
│   └── ui/               # Bubble Tea UI components
├── proto/            # Protocol Buffer definitions
├── openspec/         # Specifications and change proposals
└── logs/             # Application logs
```

## Testing

Run all tests:

```bash
go test ./... -v
```

Run specific test suites:

```bash
go test ./internal/infra/binance/... -v      # Binance client tests
go test ./internal/ai/openrouter/... -v      # OpenRouter client tests
go test ./internal/domain/indicators/... -v  # Indicator tests
```

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `OPENROUTER_API_KEY` | API key for OpenRouter |
| `PORT` | Server port (default: 50051) |
| `SYMBOL` | Default trading symbol |

### Logs

Logs are written to `logs/quantacode.log` in JSON format.

## Supported Trading Pairs

- BTCUSDT, ETHUSDT, BNBUSDT, XRPUSDT, ADAUSDT
- DOGEUSDT, SOLUSDT, DOTUSDT, MATICUSDT, LTCUSDT
- AVAXUSDT, LINKUSDT, ATOMUSDT, UNIUSDT, XLMUSDT

## License

MIT
