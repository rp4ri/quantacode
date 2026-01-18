package binance

import (
	"testing"
	"time"
)

func TestParseTicker(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPrice float64
		wantVol   float64
		wantSym   string
		wantErr   bool
	}{
		{
			name:      "binance.us real format",
			input:     `{"e":"24hrTicker","E":1768723460936,"s":"BTCUSDT","p":"-100.12000000","P":"-0.105","w":"95305.95535979","x":"95086.81000000","c":"95038.26000000","Q":"0.00021000","b":"95034.50000000","B":"0.05263000","a":"95034.63000000","A":"0.04052000","o":"95138.38000000","h":"95642.81000000","l":"94800.04000000","v":"16.80704000","q":"1601811.00397020","O":1768637060935,"C":1768723460935,"F":31019297,"L":31020323,"n":1027}`,
			wantPrice: 95038.26,
			wantVol:   16.80704,
			wantSym:   "BTCUSDT",
			wantErr:   false,
		},
		{
			name:      "binance.com format",
			input:     `{"e":"24hrTicker","E":1705574400000,"s":"ETHUSDT","c":"2500.00","v":"500.00"}`,
			wantPrice: 2500.00,
			wantVol:   500.00,
			wantSym:   "ETHUSDT",
			wantErr:   false,
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "wrong event type",
			input:   `{"e":"trade","E":1705574400000,"s":"BTCUSDT","c":"50000.00","v":"100.00"}`,
			wantErr: true,
		},
		{
			name:      "zero timestamp uses current time",
			input:     `{"e":"24hrTicker","E":0,"s":"XRPUSDT","c":"0.50","v":"1000.00"}`,
			wantPrice: 0.50,
			wantVol:   1000.00,
			wantSym:   "XRPUSDT",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTicker([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Price != tt.wantPrice {
				t.Errorf("parseTicker() Price = %v, want %v", got.Price, tt.wantPrice)
			}
			if got.Volume != tt.wantVol {
				t.Errorf("parseTicker() Volume = %v, want %v", got.Volume, tt.wantVol)
			}
			if got.Symbol != tt.wantSym {
				t.Errorf("parseTicker() Symbol = %v, want %v", got.Symbol, tt.wantSym)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("btcusdt")
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.symbol != "btcusdt" {
		t.Errorf("NewClient() symbol = %v, want btcusdt", client.symbol)
	}
	if client.simulate {
		t.Error("NewClient() should not be simulated")
	}
}

func TestNewSimulatedClient(t *testing.T) {
	client := NewSimulatedClient("ethusdt")
	if client == nil {
		t.Fatal("NewSimulatedClient() returned nil")
	}
	if client.symbol != "ethusdt" {
		t.Errorf("NewSimulatedClient() symbol = %v, want ethusdt", client.symbol)
	}
	if !client.simulate {
		t.Error("NewSimulatedClient() should be simulated")
	}
}

func TestSubscribe(t *testing.T) {
	client := NewSimulatedClient("btcusdt")
	ch := client.Subscribe()
	if ch == nil {
		t.Fatal("Subscribe() returned nil channel")
	}

	// Verify we can subscribe multiple times
	ch2 := client.Subscribe()
	if ch2 == nil {
		t.Fatal("Second Subscribe() returned nil channel")
	}

	if len(client.subscribers) != 2 {
		t.Errorf("Expected 2 subscribers, got %d", len(client.subscribers))
	}
}

func TestBroadcast(t *testing.T) {
	client := NewSimulatedClient("btcusdt")
	ch1 := client.Subscribe()
	ch2 := client.Subscribe()

	update := PriceUpdate{
		Symbol:    "BTCUSDT",
		Price:     50000.0,
		Volume:    100.0,
		Timestamp: time.Now(),
	}

	client.broadcast(update)

	select {
	case got := <-ch1:
		if got.Price != update.Price {
			t.Errorf("ch1 got Price = %v, want %v", got.Price, update.Price)
		}
	default:
		t.Error("ch1 did not receive broadcast")
	}

	select {
	case got := <-ch2:
		if got.Price != update.Price {
			t.Errorf("ch2 got Price = %v, want %v", got.Price, update.Price)
		}
	default:
		t.Error("ch2 did not receive broadcast")
	}
}
