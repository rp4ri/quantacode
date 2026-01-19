package binance

import (
	"testing"
	"time"
)

func TestParseCombinedStream(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPrice float64
		wantVol   float64
		wantSym   string
		wantErr   bool
	}{
		{
			name:      "miniTicker stream",
			input:     `{"stream":"btcusdt@miniTicker","data":{"e":"24hrMiniTicker","E":1768723460936,"s":"btcusdt","c":"95038.26","o":"95138.38","h":"95642.81","l":"94800.04","v":"16.80704","q":"1601811.00"}}`,
			wantPrice: 95038.26,
			wantVol:   16.80704,
			wantSym:   "BTCUSDT",
			wantErr:   false,
		},
		{
			name:      "aggTrade stream",
			input:     `{"stream":"ethusdt@aggTrade","data":{"e":"aggTrade","E":1705574400000,"s":"ethusdt","a":12345,"p":"2500.00","q":"0.5","f":100,"l":105,"T":1705574400000,"m":true}}`,
			wantPrice: 2500.00,
			wantVol:   0.5,
			wantSym:   "ETHUSDT",
			wantErr:   false,
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "unknown stream type",
			input:   `{"stream":"btcusdt@unknown","data":{}}`,
			wantErr: true,
		},
		{
			name:      "zero timestamp uses current time",
			input:     `{"stream":"xrpusdt@miniTicker","data":{"e":"24hrMiniTicker","E":0,"s":"xrpusdt","c":"0.50","o":"0.49","h":"0.51","l":"0.48","v":"1000.00","q":"500.00"}}`,
			wantPrice: 0.50,
			wantVol:   1000.00,
			wantSym:   "XRPUSDT",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCombinedStream([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCombinedStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Price != tt.wantPrice {
				t.Errorf("parseCombinedStream() Price = %v, want %v", got.Price, tt.wantPrice)
			}
			if got.Volume != tt.wantVol {
				t.Errorf("parseCombinedStream() Volume = %v, want %v", got.Volume, tt.wantVol)
			}
			if got.Symbol != tt.wantSym {
				t.Errorf("parseCombinedStream() Symbol = %v, want %v", got.Symbol, tt.wantSym)
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
