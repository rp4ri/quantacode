# Change: Add CLI User Interface

## Why
Implement a terminal user interface (TUI) for the QuantaCode CLI that provides real-time price updates, conversational AI chat, and technical indicator display using Bubble Tea for interactive terminal UI and Cobra for command-line parsing.

## What Changes
- Create cmd/cli/main.go with Cobra framework (root command, 'chat' subcommand)
- Add --server flag for gRPC server connection (default: localhost:50051)
- Implement Bubble Tea TUI with split-screen layout
- Create internal/ui/chat/view.go: price ticker top, chat messages bottom
- Use Lipgloss styling: green/red for price changes
- Add user input textarea at bottom of chat
- Implement gRPC client to send questions and display AI responses
- Show typing indicator while AI generates response
- Create internal/ui/indicators/panel.go: right sidebar with RSI, SMA, EMA
- Color-code RSI (red >70, blue <30, white else)
- Update indicators in real-time via gRPC streaming

## Impact
- Affected specs: cli-framework, tui-chat, tui-indicators
- New files: cmd/cli/main.go, internal/ui/chat/view.go, internal/ui/indicators/panel.go
- New dependencies: cobra, bubbletea, lipgloss, bubbletea/textinput
- New capability: Interactive TUI for real-time trading analysis
