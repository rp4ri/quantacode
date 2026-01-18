package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PriceUpdate represents a price tick from Binance.
type PriceUpdate struct {
	Symbol    string
	Price     float64
	Volume    float64
	Timestamp time.Time
}

// Client manages WebSocket connection to Binance.
type Client struct {
	symbol      string
	conn        *websocket.Conn
	mu          sync.Mutex
	subscribers []chan PriceUpdate
	subMu       sync.RWMutex
	done        chan struct{}
	simulate    bool
}

// NewClient creates a new Binance WebSocket client.
func NewClient(symbol string) *Client {
	return &Client{
		symbol: symbol,
		done:   make(chan struct{}),
	}
}

// NewSimulatedClient creates a client that generates fake price data.
func NewSimulatedClient(symbol string) *Client {
	return &Client{
		symbol:   symbol,
		done:     make(chan struct{}),
		simulate: true,
	}
}

// Subscribe adds a subscriber channel for price updates.
func (c *Client) Subscribe() <-chan PriceUpdate {
	ch := make(chan PriceUpdate, 100)
	c.subMu.Lock()
	c.subscribers = append(c.subscribers, ch)
	c.subMu.Unlock()
	return ch
}

// Connect establishes WebSocket connection and starts reading.
func (c *Client) Connect(ctx context.Context) error {
	if c.simulate {
		go c.simulateLoop(ctx)
		return nil
	}

	wsURL := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   fmt.Sprintf("/ws/%s@ticker", c.symbol),
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	header := make(map[string][]string)
	header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}

	var err error
	c.conn, _, err = dialer.DialContext(ctx, wsURL.String(), header)
	if err != nil {
		return fmt.Errorf("connect to binance: %w", err)
	}

	go c.readLoop(ctx)
	return nil
}

func (c *Client) simulateLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	price := 48000.0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			delta := (rng.Float64()*2 - 1) * 50
			price = math.Max(1, price+delta)

			update := PriceUpdate{
				Symbol:    strings.ToUpper(c.symbol),
				Price:     price,
				Volume:    rng.Float64() * 1000,
				Timestamp: time.Now(),
			}
			c.broadcast(update)
		}
	}
}

func (c *Client) readLoop(ctx context.Context) {
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			log.Printf("binance read error: %v, reconnecting...", err)
			c.reconnect(ctx)
			continue
		}

		update, err := parseTicker(message)
		if err != nil {
			log.Printf("parse ticker error: %v", err)
			continue
		}

		c.broadcast(update)
	}
}

func (c *Client) broadcast(update PriceUpdate) {
	c.subMu.RLock()
	defer c.subMu.RUnlock()

	for _, ch := range c.subscribers {
		select {
		case ch <- update:
		default:
			// drop if channel full
		}
	}
}

func (c *Client) reconnect(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}

	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		wsURL := url.URL{
			Scheme: "wss",
			Host:   "stream.binance.com:9443",
			Path:   fmt.Sprintf("/ws/%s@ticker", c.symbol),
		}

		var err error
		c.conn, _, err = websocket.DefaultDialer.DialContext(ctx, wsURL.String(), nil)
		if err == nil {
			log.Println("binance reconnected")
			return
		}

		log.Printf("reconnect failed: %v, retrying in %v", err, backoff)
		time.Sleep(backoff)
		backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
	}
}

// Close shuts down the client.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.done:
	default:
		close(c.done)
	}

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type tickerMessage struct {
	Symbol string `json:"s"`
	Price  string `json:"c"`
	Volume string `json:"v"`
	Time   int64  `json:"E"`
}

func parseTicker(data []byte) (PriceUpdate, error) {
	var msg tickerMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return PriceUpdate{}, err
	}

	var price, volume float64
	fmt.Sscanf(msg.Price, "%f", &price)
	fmt.Sscanf(msg.Volume, "%f", &volume)

	return PriceUpdate{
		Symbol:    msg.Symbol,
		Price:     price,
		Volume:    volume,
		Timestamp: time.UnixMilli(msg.Time),
	}, nil
}
