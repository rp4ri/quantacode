## Context
The QuantaCode CLI requires an interactive terminal user interface (TUI) that displays real-time price updates, allows conversational AI interaction, and shows technical indicators. Bubble Tea provides an Elm-inspired architecture for building TUIs in Go, while Lipgloss enables terminal styling. The UI must be responsive, handle streaming gRPC data, and provide a smooth user experience.

## Goals / Non-Goals
- Goals: Interactive TUI with real-time updates, conversational AI chat, technical indicator display, keyboard-driven interaction
- Non-Goals: Mouse support, graphs/charts visualization, configuration UI, multiple symbol monitoring in this phase

## Decisions

### TUI Framework: Bubble Tea
- **Decision**: Use Bubble Tea for TUI framework
- **Rationale**: Elm architecture (Model-View-Update), excellent documentation, composable components, async message handling
- **Alternative**: Rich (Python) - rejected as wrong language, termui - rejected as less maintained

### Command Framework: Cobra
- **Decision**: Use Cobra for command-line parsing
- **Rationale**: Industry standard for Go CLIs, excellent subcommand support, flag management
- **Alternative**: Urfave/cli - rejected as less feature-rich

### Layout: Split Screen with Right Sidebar
- **Decision**: Three-column layout: main area (ticker + chat + input), right sidebar (indicators)
- **Rationale**: Maximizes chat space while keeping indicators visible, follows modern terminal UI patterns
- **Alternative**: Vertical stacking - rejected for poor use of horizontal space
- **Alternative**: Popup indicators - rejected as less discoverable

### Price Color Coding: Green/Red
- **Decision**: Use green for price increase, red for decrease
- **Rationale**: Standard trading convention, intuitive for traders
- **Trade-off**: Doesn't account for green/red colorblindness (could add pattern indicator in future)

### RSI Color Coding: Red >70, Blue <30, White Else
- **Decision**: Overbought (RSI >70) red, oversold (RSI <30) blue, neutral white
- **Rationale**: Visual cues for trading signals, red/blue contrast in terminal
- **Alternative**: Green/yellow/red - rejected as less clear for terminal colors

### Typing Indicator: Animated Dots
- **Decision**: Show "AI is thinking..." with animated dots (., .., ...)
- **Rationale**: Clear feedback that AI is generating response, minimal implementation
- **Alternative**: Progress bar - rejected for streaming responses
- **Alternative**: No indicator - rejected for poor UX

### Input Method: Textarea with Enter to Send
- **Decision**: Textarea at bottom, Enter sends question, Ctrl+C or 'q' quits
- **Rationale**: Familiar chat interface, minimal key presses
- **Alternative**: Vi/vim keybindings - rejected for steep learning curve
- **Alternative**: Multi-line input - rejected for complexity

### gRPC Streaming Integration
- **Decision**: Subscribe to StreamPrices, call StreamAIAnalysis per question
- **Rationale**: Real-time updates, efficient bandwidth usage, matches server design
- **Alternative**: Polling - rejected as inefficient
- **Alternative**: Unidirectional streaming - rejected for chat interaction

## Layout Architecture

```
┌─────────────────────────────────┬─────────────────┐
│  BTC/USD: $50,123.45 (+1.2%)   │  RSI: 75.2      │
│  (green/red based on change)   │  (red)          │
├─────────────────────────────────┼─────────────────┤
│                                 │                 │
│  [User] Why is RSI so high?     │  SMA14: $50,100 │
│                                 │                 │
│  [AI] RSI >70 indicates         │  EMA12: $50,080 │
│  overbought condition...        │                 │
│                                 │                 │
│  [User] Should I sell?          │                 │
│                                 │                 │
│  [AI] Not financial advice...   │                 │
│                                 │                 │
└─────────────────────────────────┴─────────────────┘
┌─────────────────────────────────────────────────────┐
│  Type your question... [Enter to send, q to quit]   │
└─────────────────────────────────────────────────────┘
```

## Data Structures

```go
// TUI Model
type Model struct {
    serverAddr    string
    grpcClient    pb.MarketDataServiceClient
    
    // Price state
    currentPrice  *pb.PriceUpdate
    previousPrice *pb.PriceUpdate
    
    // Chat state
    chatMessages  []ChatMessage
    input         textinput.Model
    waitingForAI  bool
    
    // Indicator state
    indicators    *pb.IndicatorValues
    
    // UI state
    viewport      viewport.Model
    ready         bool
}

type ChatMessage struct {
    Role      string // "user" or "ai"
    Content   string
    Timestamp time.Time
}

// Bubble Tea Messages
type PriceUpdateMsg *pb.PriceUpdate
type ChatResponseMsg *pb.AIAnalysis
type TypingIndicatorMsg struct {
    Active bool
}

// Indicators Panel Model
type IndicatorsModel struct {
    RSI   float64
    SMA14 float64
    EMA12 float64
}
```

## Lipgloss Styling

```go
var (
    // Colors
    greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // Green
    redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))   // Red
    blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))   // Blue
    whiteStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))  // White
    
    // Layout
    tickerStyle    = lipgloss.NewStyle().Padding(1)
    chatStyle      = lipgloss.NewStyle().Padding(1)
    inputStyle     = lipgloss.NewStyle().Padding(1)
    sidebarStyle   = lipgloss.NewStyle().Padding(1).Border(lipgloss.Border{})
)
```

## Update Loop Architecture

```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    // Price updates from gRPC
    case PriceUpdateMsg:
        m.previousPrice = m.currentPrice
        m.currentPrice = msg
        return m, nil
    
    // Chat response from gRPC
    case ChatResponseMsg:
        m.chatMessages = append(m.chatMessages, ChatMessage{
            Role:    "ai",
            Content: msg.Analysis,
        })
        m.waitingForAI = false
        return m, nil
    
    // User input
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            // Send question to gRPC
            return m, sendQuestion(m.input.Value())
        case tea.KeyCtrlC, tea.KeyRunes:
            if string(msg.Runes) == "q" {
                return m, tea.Quit
            }
        }
    }
    
    // Delegate to components
    m.input = m.input.Update(msg)
    return m, nil
}
```

## Risks / Trade-offs

### Terminal Compatibility
- **Risk**: Different terminal emulators render colors differently
- **Mitigation**: Use basic ANSI colors (16-color palette), test on common terminals
- **Trade-off**: Limited color palette for compatibility vs. rich styling

### Window Size
- **Risk**: Small terminal windows break layout
- **Mitigation**: Minimum terminal size check, graceful degradation on small screens
- **Trade-off**: Reject small terminals vs. adaptive layout complexity

### Performance
- **Risk**: High-frequency price updates (100ms) cause UI lag
- **Mitigation**: Debounce updates, only render on tea.WindowSizeMsg, efficient string formatting
- **Trade-off**: Slightly delayed display vs. smooth UI

### Chat History
- **Risk**: Unlimited chat history grows memory
- **Mitigation**: Limit to last 100 messages, implement pagination if needed
- **Trade-off**: Limited history vs. memory usage

### gRPC Connection Failures
- **Risk**: Server connection failure renders UI unusable
- **Mitigation**: Show clear error message, allow reconnection, fallback mode (if possible)
- **Trade-off**: Fail fast vs. degraded functionality

## Message Flow

```
User Input (Enter)
    ↓
CLI → gRPC StreamAIAnalysis(question)
    ↓
Server → OpenRouter API
    ↓
Server → gRPC Stream (AIAnalysis messages)
    ↓
CLI Update Loop (ChatResponseMsg)
    ↓
View renders AI response in chat area
```

```
Binance → Server → gRPC StreamPrices
    ↓
CLI Update Loop (PriceUpdateMsg)
    ↓
View renders price ticker (green/red)
    ↓
Update indicator values
    ↓
View renders indicator panel (RSI color-coded)
```

## Error Handling Strategy

### gRPC Connection Errors
- Connection refused: Show error message, offer to reconnect or exit
- Connection lost during session: Show notification, attempt reconnection
- Timeout: Show "server slow, still trying..." message

### User Input Errors
- Empty input: Ignore, show subtle feedback
- Very long input: Truncate with indicator, warn user

### Rendering Errors
- Terminal too small: Show error message with minimum size requirements
- Color not supported: Fallback to monochrome styling

## Migration Plan
N/A - new capability, no existing code to migrate

## Open Questions
- Should we support dark/light terminal background themes?
- Should we add keyboard shortcuts for common actions (clear chat, save conversation)?
- Should we display the symbol in the title bar?
- Should we support multi-line input for complex questions?
- Should we add a command palette for actions?
