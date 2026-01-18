package indicators

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/lipgloss"
    domainindicators "github.com/yourusername/quantacode/internal/domain/indicators"
)

var (
    subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
    highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    dimText   = lipgloss.Color("#666666")
    greenColor = lipgloss.Color("#73F59F")
    redColor   = lipgloss.Color("#FF6B6B")
    blueColor  = lipgloss.Color("#6CB6FF")
    purpleColor = lipgloss.Color("#D4A5FF")
)

// Panel renders indicator values in a right sidebar.
type Panel struct {
    width   int
    height  int
    history domainindicators.IndicatorHistory
}

// NewPanel creates a Panel with a default width.
func NewPanel() Panel {
    return Panel{width: 32, height: 30}
}

// WithWidth updates the target width.
func (p Panel) WithWidth(width int) Panel {
    if width < 28 {
        width = 28
    }
    p.width = width
    return p
}

// WithHeight updates the target height.
func (p Panel) WithHeight(height int) Panel {
    p.height = height
    return p
}

// WithHistory updates the indicator history.
func (p Panel) WithHistory(history domainindicators.IndicatorHistory) Panel {
    p.history = history
    return p
}

// View renders the indicator state.
func (p Panel) View(vals domainindicators.AggregatedValues) string {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(highlight).
        Render("◆ Indicadores")
    
    border := lipgloss.NewStyle().
        Foreground(subtle).
        Render(strings.Repeat("─", p.width-2))

    currentSection := p.renderCurrentValues(vals)

    content := lipgloss.JoinVertical(lipgloss.Left,
        title,
        border,
        currentSection,
    )

    return lipgloss.NewStyle().
        Width(p.width).
        Padding(0, 1).
        Render(content)
}

func (p Panel) renderCurrentValues(vals domainindicators.AggregatedValues) string {
    labelStyle := lipgloss.NewStyle().Foreground(dimText)
    warmingStyle := lipgloss.NewStyle().Foreground(dimText).Italic(true)
    
    // Check if indicators are still warming up (RSI=0 and SMA=0 means not enough data)
    isWarmingUp := vals.RSI == 0 && vals.SMA == 0
    
    var rsiLine, smaLine, emaLine string
    
    if isWarmingUp {
        rsiLine = fmt.Sprintf("%s %s", 
            labelStyle.Render("RSI:"),
            warmingStyle.Render("calentando..."))
        smaLine = fmt.Sprintf("%s %s",
            labelStyle.Render("SMA:"),
            warmingStyle.Render("calentando..."))
        emaLine = fmt.Sprintf("%s %s",
            labelStyle.Render("EMA:"),
            warmingStyle.Render("calentando..."))
    } else {
        rsiStyle := lipgloss.NewStyle().Bold(true)
        rsiLabel := "RSI"
        switch {
        case vals.RSI >= 70:
            rsiStyle = rsiStyle.Foreground(redColor)
            rsiLabel = "RSI ⚠"
        case vals.RSI <= 30 && vals.RSI > 0:
            rsiStyle = rsiStyle.Foreground(blueColor)
            rsiLabel = "RSI ⚠"
        default:
            rsiStyle = rsiStyle.Foreground(lipgloss.Color("#FFFFFF"))
        }

        rsiLine = fmt.Sprintf("%s %s", 
            labelStyle.Render(rsiLabel+":"),
            rsiStyle.Render(fmt.Sprintf("%.2f", vals.RSI)))
        
        smaLine = fmt.Sprintf("%s %s",
            labelStyle.Render("SMA:"),
            lipgloss.NewStyle().Foreground(greenColor).Render(fmt.Sprintf("%.2f", vals.SMA)))
        
        emaLine = fmt.Sprintf("%s %s",
            labelStyle.Render("EMA:"),
            lipgloss.NewStyle().Foreground(purpleColor).Render(fmt.Sprintf("%.2f", vals.EMA)))
    }

    return lipgloss.JoinVertical(lipgloss.Left, rsiLine, smaLine, emaLine)
}

func (p Panel) renderHistory() string {
    if len(p.history.RSI) == 0 {
        return lipgloss.NewStyle().
            Foreground(dimText).
            Italic(true).
            Render("Esperando datos...")
    }

    headerStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(dimText)
    
    header := headerStyle.Render(fmt.Sprintf("Historial (%d velas)", len(p.history.RSI)))
    
    tableHeader := lipgloss.NewStyle().
        Foreground(dimText).
        Render(" #   RSI    SMA      EMA")
    
    var rows []string
    rows = append(rows, header)
    rows = append(rows, tableHeader)
    
    displayCount := len(p.history.RSI)
    if displayCount > 15 {
        displayCount = 15
    }
    
    startIdx := len(p.history.RSI) - displayCount
    
    for i := startIdx; i < len(p.history.RSI); i++ {
        rowNum := i - startIdx + 1
        rsi := p.history.RSI[i]
        sma := p.history.SMA[i]
        ema := p.history.EMA[i]
        
        rsiStyle := lipgloss.NewStyle()
        if rsi >= 70 {
            rsiStyle = rsiStyle.Foreground(redColor)
        } else if rsi <= 30 && rsi != 0 {
            rsiStyle = rsiStyle.Foreground(blueColor)
        } else {
            rsiStyle = rsiStyle.Foreground(lipgloss.Color("#888888"))
        }
        
        row := fmt.Sprintf("%2d  %s %7.1f %7.1f",
            rowNum,
            rsiStyle.Render(fmt.Sprintf("%5.1f", rsi)),
            sma,
            ema,
        )
        rows = append(rows, row)
    }
    
    if len(p.history.RSI) > 15 {
        moreMsg := lipgloss.NewStyle().
            Foreground(dimText).
            Italic(true).
            Render(fmt.Sprintf("  ... +%d más", len(p.history.RSI)-15))
        rows = append(rows, moreMsg)
    }

    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
