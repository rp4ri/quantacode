package openrouter

import (
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("NewClient() apiKey = %v, want test-api-key", client.apiKey)
	}
	if client.model != defaultModel {
		t.Errorf("NewClient() model = %v, want %v", client.model, defaultModel)
	}
}

func TestWithModel(t *testing.T) {
	client := NewClient("test-key").WithModel("custom/model")
	if client.model != "custom/model" {
		t.Errorf("WithModel() model = %v, want custom/model", client.model)
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		name        string
		symbol      string
		price       float64
		rsi         float64
		sma         float64
		ema         float64
		history     *IndicatorHistory
		wantContain []string
	}{
		{
			name:    "basic prompt without history",
			symbol:  "BTCUSDT",
			price:   50000.0,
			rsi:     55.5,
			sma:     49000.0,
			ema:     49500.0,
			history: nil,
			wantContain: []string{
				"BTCUSDT",
				"$50000.00",
				"55.50",
				"49000.00",
				"49500.00",
			},
		},
		{
			name:   "prompt with history",
			symbol: "ETHUSDT",
			price:  2500.0,
			rsi:    70.0,
			sma:    2400.0,
			ema:    2450.0,
			history: &IndicatorHistory{
				RSI: []float64{65.0, 68.0, 70.0},
				SMA: []float64{2350.0, 2380.0, 2400.0},
				EMA: []float64{2400.0, 2430.0, 2450.0},
			},
			wantContain: []string{
				"ETHUSDT",
				"Historial de indicadores",
				"65.00",
				"68.00",
				"70.00",
			},
		},
		{
			name:   "prompt with empty history",
			symbol: "XRPUSDT",
			price:  0.5,
			rsi:    30.0,
			sma:    0.48,
			ema:    0.49,
			history: &IndicatorHistory{
				RSI: []float64{},
				SMA: []float64{},
				EMA: []float64{},
			},
			wantContain: []string{
				"XRPUSDT",
				"$0.50",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.buildSystemPrompt(tt.symbol, tt.price, tt.rsi, tt.sma, tt.ema, tt.history)
			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("buildSystemPrompt() missing %q in output:\n%s", want, got)
				}
			}
		})
	}
}

func TestBuildSystemPromptNoAutoAnalysis(t *testing.T) {
	client := NewClient("test-key")
	prompt := client.buildSystemPrompt("BTCUSDT", 50000, 50, 49000, 49500, nil)

	// Verify the prompt instructs AI not to auto-analyze
	mustContain := []string{
		"Solo proporciona análisis técnico cuando el usuario lo solicite",
		"Si el usuario hace una pregunta general o saluda",
	}

	for _, want := range mustContain {
		if !strings.Contains(prompt, want) {
			t.Errorf("buildSystemPrompt() should contain instruction %q", want)
		}
	}
}

func TestIndicatorHistoryStruct(t *testing.T) {
	history := &IndicatorHistory{
		RSI: []float64{50.0, 55.0, 60.0},
		SMA: []float64{100.0, 101.0, 102.0},
		EMA: []float64{99.0, 100.0, 101.0},
	}

	if len(history.RSI) != 3 {
		t.Errorf("IndicatorHistory RSI length = %d, want 3", len(history.RSI))
	}
	if history.RSI[0] != 50.0 {
		t.Errorf("IndicatorHistory RSI[0] = %v, want 50.0", history.RSI[0])
	}
}
