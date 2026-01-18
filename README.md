# QuantaCode CLI

QuantaCode is an AI-assisted trading analysis tool that streams live market data, computes technical indicators (RSI, SMA, EMA), and renders an interactive terminal UI powered by Bubble Tea.

## Prerequisites

- Go 1.21+
- `protoc` (Protocol Buffers compiler)
- Binance internet access for real-time data

Install protoc on Ubuntu:

```bash
sudo apt update && sudo apt install -y protobuf-compiler
```

## Install Dependencies

```bash
go mod tidy
```

## Running the CLI

```bash
go run ./cmd/cli chat --server=localhost:50051 --symbol=BTCUSDT --openrouter-key=$OPENROUTER_KEY
```

Flags:

- `--server` (default `localhost:50051`): gRPC server address
- `--symbol` (default `BTCUSDT`): Trading pair
- `--openrouter-key`: API key for AI analysis (can also come from `OPENROUTER_KEY` env var)

## Features (in progress)

- Cobra-based CLI entrypoint (`cmd/cli/main.go`)
- Bubble Tea TUI with ticker, chat, indicators panel
- Typing indicator animations
- Live gRPC streaming for price updates and AI responses

## Development Workflow

1. Start the gRPC server (see `cmd/server/main.go` once implemented).
2. Run the CLI as shown above.
3. Use `q` or `Ctrl+C` to exit the TUI.

## Testing

All tests live under dedicated `tests/` folders. Example:

```bash
go test ./internal/domain/indicators/tests -v
```

More test suites will be added as CLI and server modules land.
