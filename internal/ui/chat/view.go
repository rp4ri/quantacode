package chat

import (
    "context"
    "fmt"
    "math"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "github.com/yourusername/quantacode/internal/ai/openrouter"
    domainindicators "github.com/yourusername/quantacode/internal/domain/indicators"
    grpcclient "github.com/yourusername/quantacode/internal/grpc/client"
    "github.com/yourusername/quantacode/internal/logging"
    indicatorpanel "github.com/yourusername/quantacode/internal/ui/indicators"
)

const (
    brandName = "QuantaCode"
)

var availablePairs = []string{
    "btcusdt", "ethusdt", "bnbusdt", "xrpusdt", "adausdt",
    "dogeusdt", "solusdt", "dotusdt", "maticusdt", "ltcusdt",
    "avaxusdt", "linkusdt", "atomusdt", "uniusdt", "xlmusdt",
}

type slashCommand struct {
    name        string
    description string
}

var slashCommands = []slashCommand{
    {name: "/clear", description: "Limpiar historial del chat"},
    {name: "/pairs", description: "Cambiar par de trading"},
}

// Config contains runtime configuration for the chat UI.
type Config struct {
    ServerAddr    string
    Symbol        string
    OpenRouterKey string
}

// Run starts the Bubble Tea program for the chat UI.
func Run(ctx context.Context, cfg Config) error {
    _ = logging.Init("logs/quantacode.log", logging.DEBUG)
    defer logging.Close()

    logger := logging.GetLogger("chat-ui")
    logger.Info("Starting chat UI")

    panel := indicatorpanel.NewPanel()
    model := newModel(cfg, panel)

    program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithContext(ctx))
    _, err := program.Run()
    return err
}

// ------ Bubble Tea model ------

type model struct {
    cfg          Config
    textarea     textarea.Model
    viewport     viewport.Model
    messages     []chatMessage
    typing       bool
    typingFrame  int
    ready        bool
    width        int
    height       int
    currentPrice float64
    priceChange  float64
    prevPrice    float64

    indicatorValues domainindicators.AggregatedValues
    panel           indicatorpanel.Panel

    grpcClient  *grpcclient.Client
    connected   bool
    priceCh     chan grpcclient.PriceUpdate
    indicatorCh chan grpcclient.IndicatorUpdate

    aiClient        *openrouter.Client
    streamingMsg    string
    aiStreamCh      <-chan openrouter.StreamChunk
    logger          *logging.Logger
    indicatorHistory *openrouter.IndicatorHistory
    
    inputHistory      []string
    historyIndex      int
    showPairSelect    bool
    pairSelectIndex   int
    showSlashMenu     bool
    slashMenuIndex    int
}

type chatMessage struct {
    author    string
    content   string
    timestamp time.Time
}

func newModel(cfg Config, panel indicatorpanel.Panel) model {
    ti := textarea.New()
    ti.Placeholder = "Escribe tu pregunta..."
    ti.Focus()
    ti.Prompt = "› "
    ti.CharLimit = 500
    ti.SetHeight(3)
    ti.ShowLineNumbers = false
    ti.FocusedStyle.CursorLine = lipgloss.NewStyle()
    ti.BlurredStyle.CursorLine = lipgloss.NewStyle()
    ti.FocusedStyle.Base = lipgloss.NewStyle()
    ti.BlurredStyle.Base = lipgloss.NewStyle()

    vp := viewport.New(80, 20)
    vp.Style = lipgloss.NewStyle()
    vp.MouseWheelEnabled = true
    vp.MouseWheelDelta = 3
    vp.YPosition = 0

    var aiClient *openrouter.Client
    if cfg.OpenRouterKey != "" {
        aiClient = openrouter.NewClient(cfg.OpenRouterKey)
    }

    return model{
        cfg:          cfg,
        textarea:     ti,
        viewport:     vp,
        panel:        panel,
        currentPrice: 0,
        aiClient:     aiClient,
        logger:       logging.GetLogger("chat-model"),
    }
}

func (m model) Init() tea.Cmd {
    return tea.Batch(connectCmd(m.cfg.ServerAddr, m.cfg.Symbol), typingTickerCmd())
}

type priceUpdateMsg struct {
    price  float64
    symbol string
}
type indicatorUpdateMsg struct {
    rsi        float64
    sma        float64
    ema        float64
    rsiHistory []float64
    smaHistory []float64
    emaHistory []float64
}
type typingTickMsg struct{}
type aiResponseMsg struct {
    content string
}
type errMsg struct {
    err error
}
type connectedMsg struct {
    client *grpcclient.Client
}

func connectCmd(serverAddr, symbol string) tea.Cmd {
    return func() tea.Msg {
        client, err := grpcclient.New(serverAddr)
        if err != nil {
            return errMsg{err: fmt.Errorf("connect: %w", err)}
        }
        return connectedMsg{client: client}
    }
}

func typingTickerCmd() tea.Cmd {
    return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
        return typingTickMsg{}
    })
}

type aiStreamChunkMsg struct {
    content  string
    done     bool
    err      error
    streamCh <-chan openrouter.StreamChunk
}

func startAIStreamCmd(client *openrouter.Client, prompt, symbol string, price, rsi, sma, ema float64, history *openrouter.IndicatorHistory) tea.Cmd {
    return func() tea.Msg {
        if client == nil {
            return aiStreamChunkMsg{content: "Error: API key no configurada", done: true}
        }

        ctx := context.Background()
        chunkCh, err := client.StreamAnalysis(ctx, prompt, symbol, price, rsi, sma, ema, history)
        if err != nil {
            return aiStreamChunkMsg{err: err, done: true}
        }

        chunk, ok := <-chunkCh
        if !ok {
            return aiStreamChunkMsg{done: true}
        }
        return aiStreamChunkMsg{content: chunk.Content, done: chunk.Done, err: chunk.Error, streamCh: chunkCh}
    }
}

func continueAIStreamCmd(chunkCh <-chan openrouter.StreamChunk) tea.Cmd {
    return func() tea.Msg {
        chunk, ok := <-chunkCh
        if !ok {
            return aiStreamChunkMsg{done: true}
        }
        return aiStreamChunkMsg{content: chunk.Content, done: chunk.Done, err: chunk.Error, streamCh: chunkCh}
    }
}

type startStreamMsg struct {
    priceCh     chan grpcclient.PriceUpdate
    indicatorCh chan grpcclient.IndicatorUpdate
}

func startStreamCmd(client *grpcclient.Client, symbol string) tea.Cmd {
    return func() tea.Msg {
        priceCh := make(chan grpcclient.PriceUpdate, 100)
        indicatorCh := make(chan grpcclient.IndicatorUpdate, 100)

        go func() {
            _ = client.StreamPrices(context.Background(), symbol, priceCh, indicatorCh)
        }()

        return startStreamMsg{priceCh: priceCh, indicatorCh: indicatorCh}
    }
}

func waitForUpdateCmd(priceCh <-chan grpcclient.PriceUpdate, indicatorCh <-chan grpcclient.IndicatorUpdate) tea.Cmd {
    return func() tea.Msg {
        select {
        case p, ok := <-priceCh:
            if !ok {
                return errMsg{err: fmt.Errorf("price channel closed")}
            }
            return priceUpdateMsg{price: p.Price, symbol: p.Symbol}
        case i, ok := <-indicatorCh:
            if !ok {
                return errMsg{err: fmt.Errorf("indicator channel closed")}
            }
            return indicatorUpdateMsg{
                rsi:        i.RSI,
                sma:        i.SMA,
                ema:        i.EMA,
                rsiHistory: i.RSIHistory,
                smaHistory: i.SMAHistory,
                emaHistory: i.EMAHistory,
            }
        }
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.ready = true
        contentWidth := m.contentWidth()
        if contentWidth > 0 {
            m.textarea.SetWidth(contentWidth - 4)
        }
        headerHeight := 3
        inputHeight := 5
        chatHeight := m.height - headerHeight - inputHeight - 2
        if chatHeight < 5 {
            chatHeight = 5
        }
        m.viewport.Width = contentWidth - 2
        m.viewport.Height = chatHeight
        m.updateViewportContent()

    case tea.KeyMsg:
        // Handle pair selection modal
        if m.showPairSelect {
            switch msg.Type {
            case tea.KeyEsc:
                m.showPairSelect = false
            case tea.KeyUp:
                if m.pairSelectIndex > 0 {
                    m.pairSelectIndex--
                }
            case tea.KeyDown:
                if m.pairSelectIndex < len(availablePairs)-1 {
                    m.pairSelectIndex++
                }
            case tea.KeyEnter:
                selectedPair := availablePairs[m.pairSelectIndex]
                oldPair := m.cfg.Symbol
                m.showPairSelect = false
                m.cfg.Symbol = selectedPair
                // Reset price and indicators for new pair
                m.currentPrice = 0
                m.prevPrice = 0
                m.priceChange = 0
                m.indicatorValues = domainindicators.AggregatedValues{}
                m.indicatorHistory = nil
                m.priceCh = nil
                m.indicatorCh = nil
                m.logger.LogPairSwitch(oldPair, selectedPair)
                m.messages = append(m.messages, chatMessage{author: "Sistema", content: fmt.Sprintf("Cambiando a par: %s", strings.ToUpper(selectedPair)), timestamp: time.Now()})
                if m.grpcClient != nil {
                    cmds = append(cmds, startStreamCmd(m.grpcClient, selectedPair))
                }
            }
            break
        }
        
        // Handle slash command menu
        if m.showSlashMenu {
            switch msg.Type {
            case tea.KeyEsc:
                m.showSlashMenu = false
            case tea.KeyUp:
                if m.slashMenuIndex > 0 {
                    m.slashMenuIndex--
                }
            case tea.KeyDown:
                if m.slashMenuIndex < len(slashCommands)-1 {
                    m.slashMenuIndex++
                }
            case tea.KeyEnter, tea.KeyTab:
                m.textarea.SetValue(slashCommands[m.slashMenuIndex].name)
                m.showSlashMenu = false
            default:
                var cmd tea.Cmd
                m.textarea, cmd = m.textarea.Update(msg)
                cmds = append(cmds, cmd)
                // Check if still typing a slash command
                val := m.textarea.Value()
                if !strings.HasPrefix(val, "/") || strings.Contains(val, " ") {
                    m.showSlashMenu = false
                }
            }
            break
        }
        
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
            return m, tea.Quit
        case tea.KeyUp:
            if len(m.inputHistory) > 0 && m.historyIndex > 0 {
                m.historyIndex--
                m.textarea.SetValue(m.inputHistory[m.historyIndex])
            }
        case tea.KeyDown:
            if len(m.inputHistory) > 0 && m.historyIndex < len(m.inputHistory)-1 {
                m.historyIndex++
                m.textarea.SetValue(m.inputHistory[m.historyIndex])
            } else if m.historyIndex >= len(m.inputHistory)-1 {
                m.historyIndex = len(m.inputHistory)
                m.textarea.SetValue("")
            }
        case tea.KeyPgUp:
            m.viewport.LineUp(5)
        case tea.KeyPgDown:
            m.viewport.LineDown(5)
        case tea.KeyEnter:
            if m.typing {
                break
            }
            input := strings.TrimSpace(strings.ReplaceAll(m.textarea.Value(), "\n", ""))
            if input == "" {
                break
            }
            
            // Handle slash commands
            if strings.HasPrefix(input, "/") {
                cmd := strings.ToLower(strings.TrimSpace(input))
                switch {
                case cmd == "/clear" || strings.HasPrefix(cmd, "/clear "):
                    m.messages = nil
                    m.textarea.Reset()
                    m.updateViewportContent()
                    break
                case cmd == "/pairs" || strings.HasPrefix(cmd, "/pairs "):
                    m.showPairSelect = true
                    m.pairSelectIndex = 0
                    for i, p := range availablePairs {
                        if p == strings.ToLower(m.cfg.Symbol) {
                            m.pairSelectIndex = i
                            break
                        }
                    }
                    m.textarea.Reset()
                    break
                default:
                    m.messages = append(m.messages, chatMessage{author: "Sistema", content: "Comandos disponibles: /clear, /pairs", timestamp: time.Now()})
                    m.textarea.Reset()
                }
                break
            }
            
            m.inputHistory = append(m.inputHistory, input)
            m.historyIndex = len(m.inputHistory)
            m.logger.LogIO("User input", input, nil, 0)
            m.messages = append(m.messages, chatMessage{author: "Tú", content: input, timestamp: time.Now()})
            m.textarea.Reset()
            m.typing = true
            m.streamingMsg = ""
            cmds = append(cmds, startAIStreamCmd(m.aiClient, input, m.cfg.Symbol, m.currentPrice, m.indicatorValues.RSI, m.indicatorValues.SMA, m.indicatorValues.EMA, m.indicatorHistory))
        default:
            var cmd tea.Cmd
            m.textarea, cmd = m.textarea.Update(msg)
            cmds = append(cmds, cmd)
            
            // Check if user just typed "/" to show slash menu
            val := m.textarea.Value()
            if val == "/" {
                m.showSlashMenu = true
                m.slashMenuIndex = 0
            } else if strings.HasPrefix(val, "/") && !strings.Contains(val, " ") && len(val) > 1 {
                // Filter commands based on input
                m.showSlashMenu = true
            } else {
                m.showSlashMenu = false
            }
        }

    case tea.MouseMsg:
        // Always forward mouse messages to viewport for scroll handling
        var vpCmd tea.Cmd
        m.viewport, vpCmd = m.viewport.Update(msg)
        cmds = append(cmds, vpCmd)

    case connectedMsg:
        m.grpcClient = msg.client
        m.connected = true
        m.messages = append(m.messages, chatMessage{author: "Sistema", content: fmt.Sprintf("Conectado a %s (par %s)", m.cfg.ServerAddr, m.cfg.Symbol), timestamp: time.Now()})
        cmds = append(cmds, startStreamCmd(m.grpcClient, m.cfg.Symbol))

    case startStreamMsg:
        m.priceCh = msg.priceCh
        m.indicatorCh = msg.indicatorCh
        cmds = append(cmds, waitForUpdateCmd(m.priceCh, m.indicatorCh))

    case priceUpdateMsg:
        m.prevPrice = m.currentPrice
        m.currentPrice = msg.price
        m.priceChange = m.currentPrice - m.prevPrice
        m.logger.LogPriceUpdate(msg.symbol, msg.price, 0)
        if m.priceCh != nil {
            cmds = append(cmds, waitForUpdateCmd(m.priceCh, m.indicatorCh))
        }

    case indicatorUpdateMsg:
        m.indicatorValues = domainindicators.AggregatedValues{
            RSI: msg.rsi,
            SMA: msg.sma,
            EMA: msg.ema,
        }
        m.logger.LogIndicatorUpdate(msg.rsi, msg.sma, msg.ema)
        history := domainindicators.IndicatorHistory{
            RSI: msg.rsiHistory,
            SMA: msg.smaHistory,
            EMA: msg.emaHistory,
        }
        m.panel = m.panel.WithHistory(history)
        m.indicatorHistory = &openrouter.IndicatorHistory{
            RSI: msg.rsiHistory,
            SMA: msg.smaHistory,
            EMA: msg.emaHistory,
        }
        if m.priceCh != nil {
            cmds = append(cmds, waitForUpdateCmd(m.priceCh, m.indicatorCh))
        }

    case typingTickMsg:
        if m.typing {
            m.typingFrame = (m.typingFrame + 1) % len(typingFrames)
        }
        cmds = append(cmds, typingTickerCmd())

    case aiStreamChunkMsg:
        if msg.err != nil {
            m.typing = false
            m.messages = append(m.messages, chatMessage{author: "Error", content: msg.err.Error(), timestamp: time.Now()})
            m.streamingMsg = ""
            m.aiStreamCh = nil
        } else if msg.done {
            m.typing = false
            m.typingFrame = 0
            if m.streamingMsg != "" {
                m.messages = append(m.messages, chatMessage{author: brandName, content: m.streamingMsg, timestamp: time.Now()})
                m.logger.LogIO("AI response complete", nil, m.streamingMsg, 0)
            }
            m.streamingMsg = ""
            m.aiStreamCh = nil
        } else {
            m.streamingMsg += msg.content
            m.aiStreamCh = msg.streamCh
            if msg.streamCh != nil {
                cmds = append(cmds, continueAIStreamCmd(msg.streamCh))
            }
        }

    case errMsg:
        m.messages = append(m.messages, chatMessage{author: "Error", content: msg.err.Error(), timestamp: time.Now()})
        m.typing = false
        m.typingFrame = 0
    }

    m.panel = m.panel.WithWidth(m.panelWidth())
    m.updateViewportContent()
    return m, tea.Batch(cmds...)
}

var typingFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var (
    subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
    highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
    dimText   = lipgloss.Color("#666666")
    userColor = lipgloss.Color("#6CB6FF")
    aiColor   = lipgloss.Color("#D4A5FF")
    errColor  = lipgloss.Color("#FF6B6B")
    sysColor  = lipgloss.Color("#73F59F")
)

func (m *model) updateViewportContent() {
    content := m.buildChatContent()
    m.viewport.SetContent(content)
    m.viewport.GotoBottom()
}

func (m model) buildChatContent() string {
    if len(m.messages) == 0 && m.streamingMsg == "" {
        return lipgloss.NewStyle().
            Foreground(dimText).
            Italic(true).
            Render("\n  Escribe una pregunta sobre el mercado para comenzar el análisis...\n")
    }

    var builder strings.Builder
    width := m.viewport.Width - 2

    for _, msg := range m.messages {
        builder.WriteString(m.formatMessage(msg, width))
        builder.WriteString("\n")
    }

    if m.streamingMsg != "" {
        streamHeader := lipgloss.NewStyle().
            Bold(true).
            Foreground(aiColor).
            Render(brandName)
        timestamp := lipgloss.NewStyle().
            Foreground(dimText).
            Render(" • " + time.Now().Format("15:04:05"))
        builder.WriteString("\n" + streamHeader + timestamp + "\n")
        
        streamContent := lipgloss.NewStyle().
            Width(width).
            Foreground(lipgloss.Color("#FFFFFF")).
            Render(m.streamingMsg)
        builder.WriteString(streamContent)
        
        if m.typing {
            cursor := lipgloss.NewStyle().
                Foreground(aiColor).
                Render("▋")
            builder.WriteString(cursor)
        }
        builder.WriteString("\n")
    }

    return builder.String()
}

func (m model) formatMessage(msg chatMessage, width int) string {
    var authorStyle lipgloss.Style
    switch msg.author {
    case brandName:
        authorStyle = lipgloss.NewStyle().Bold(true).Foreground(aiColor)
    case "Tú":
        authorStyle = lipgloss.NewStyle().Bold(true).Foreground(userColor)
    case "Error":
        authorStyle = lipgloss.NewStyle().Bold(true).Foreground(errColor)
    case "Sistema":
        authorStyle = lipgloss.NewStyle().Bold(true).Foreground(sysColor)
    default:
        authorStyle = lipgloss.NewStyle().Bold(true).Foreground(userColor)
    }

    header := authorStyle.Render(msg.author)
    timestamp := lipgloss.NewStyle().Foreground(dimText).Render(" • " + msg.timestamp.Format("15:04:05"))
    
    content := lipgloss.NewStyle().
        Width(width).
        Render(msg.content)

    return "\n" + header + timestamp + "\n" + content
}

func (m model) View() string {
    if !m.ready {
        return lipgloss.NewStyle().
            Foreground(dimText).
            Render(" Cargando " + brandName + "...")
    }

    if m.showPairSelect {
        return m.renderPairSelect()
    }

    contentWidth := m.contentWidth()
    
    headerView := m.renderHeader(contentWidth)
    chatView := m.viewport.View()
    typingView := m.renderTyping()
    inputView := m.renderInput(contentWidth)

    leftContent := lipgloss.JoinVertical(lipgloss.Left,
        headerView,
        chatView,
        typingView,
        inputView,
    )
    
    leftBox := lipgloss.NewStyle().
        Width(contentWidth).
        Height(m.height).
        Render(leftContent)
    
    // Only show indicator panel if terminal is wide enough
    if m.panelWidth() == 0 {
        return leftBox
    }
    
    panelView := m.panel.View(m.indicatorValues)

    return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, panelView)
}

func (m model) renderPairSelect() string {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(highlight).
        Render("◆ Seleccionar Par de Trading")
    
    subtitle := lipgloss.NewStyle().
        Foreground(dimText).
        Render("Usa ↑/↓ para navegar, Enter para seleccionar, Esc para cancelar\n")
    
    var items strings.Builder
    for i, pair := range availablePairs {
        cursor := "  "
        style := lipgloss.NewStyle().Foreground(dimText)
        if i == m.pairSelectIndex {
            cursor = "› "
            style = lipgloss.NewStyle().Bold(true).Foreground(highlight)
        }
        if strings.ToLower(m.cfg.Symbol) == pair {
            items.WriteString(cursor + style.Render(strings.ToUpper(pair)+" (actual)") + "\n")
        } else {
            items.WriteString(cursor + style.Render(strings.ToUpper(pair)) + "\n")
        }
    }
    
    box := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(highlight).
        Padding(1, 2).
        Width(40).
        Render(lipgloss.JoinVertical(lipgloss.Left, title, subtitle, items.String()))
    
    return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m model) renderHeader(width int) string {
    logo := lipgloss.NewStyle().
        Bold(true).
        Foreground(highlight).
        Render("◆ " + brandName)
    
    var statusIcon, statusText string
    if m.connected {
        statusIcon = lipgloss.NewStyle().Foreground(sysColor).Render("●")
        statusText = lipgloss.NewStyle().Foreground(dimText).Render(" conectado")
    } else {
        statusIcon = lipgloss.NewStyle().Foreground(errColor).Render("○")
        statusText = lipgloss.NewStyle().Foreground(dimText).Render(" desconectado")
    }
    
    symbol := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FFFFFF")).
        Render(m.cfg.Symbol)
    
    priceStyle := lipgloss.NewStyle().Bold(true)
    changeStr := ""
    if m.priceChange > 0 {
        priceStyle = priceStyle.Foreground(sysColor)
        changeStr = fmt.Sprintf(" ↑%.2f", m.priceChange)
    } else if m.priceChange < 0 {
        priceStyle = priceStyle.Foreground(errColor)
        changeStr = fmt.Sprintf(" ↓%.2f", math.Abs(m.priceChange))
    }
    price := priceStyle.Render(fmt.Sprintf("$%.2f%s", m.currentPrice, changeStr))
    
    left := logo + "  " + statusIcon + statusText
    right := symbol + " " + price
    
    gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 2
    if gap < 1 {
        gap = 1
    }
    
    header := left + strings.Repeat(" ", gap) + right
    
    border := lipgloss.NewStyle().
        Foreground(subtle).
        Render(strings.Repeat("─", width))
    
    return header + "\n" + border
}

func (m model) renderInput(width int) string {
    border := lipgloss.NewStyle().
        Foreground(subtle).
        Render(strings.Repeat("─", width))
    
    var slashMenu string
    if m.showSlashMenu {
        slashMenu = m.renderSlashMenu()
    }
    
    prompt := lipgloss.NewStyle().
        Foreground(highlight).
        Render("› ")
    
    input := m.textarea.View()
    
    return border + "\n" + slashMenu + prompt + input
}

func (m model) renderSlashMenu() string {
    var items strings.Builder
    inputVal := strings.ToLower(m.textarea.Value())
    
    for i, cmd := range slashCommands {
        // Filter based on what user typed
        if inputVal != "/" && !strings.HasPrefix(cmd.name, inputVal) {
            continue
        }
        
        cursor := "  "
        nameStyle := lipgloss.NewStyle().Foreground(dimText)
        descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
        
        if i == m.slashMenuIndex {
            cursor = "› "
            nameStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
            descStyle = lipgloss.NewStyle().Foreground(dimText)
        }
        
        items.WriteString(cursor + nameStyle.Render(cmd.name) + " " + descStyle.Render(cmd.description) + "\n")
    }
    
    if items.Len() == 0 {
        return ""
    }
    
    box := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(subtle).
        Padding(0, 1).
        Render(items.String())
    
    return box + "\n"
}

func (m model) renderTyping() string {
    if !m.typing || m.streamingMsg != "" {
        return ""
    }
    spinner := typingFrames[m.typingFrame%len(typingFrames)]
    return lipgloss.NewStyle().
        Foreground(aiColor).
        Render("  " + spinner + " " + brandName + " está pensando...")
}

func (m model) contentWidth() int {
    panelWidth := m.panelWidth()
    contentWidth := m.width - panelWidth - 1
    if contentWidth < 40 {
        contentWidth = 40
    }
    return contentWidth
}

func (m model) panelWidth() int {
    // Hide panel if terminal too narrow (< 100 columns)
    if m.width < 100 {
        return 0
    }
    if m.width == 0 {
        return 32
    }
    w := m.width / 4
    if w < 32 {
        w = 32
    }
    if w > 50 {
        w = 50
    }
    return w
}
