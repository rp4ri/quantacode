package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Kline represents a single candlestick from Binance.
type Kline struct {
	OpenTime  time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime time.Time
}

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
// Tries binance.us first (for US servers), then falls back to binance.com
func (c *Client) Connect(ctx context.Context) error {
	if c.simulate {
		go c.simulateLoop(ctx)
		return nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	header := make(map[string][]string)
	header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}

	// Try binance.com first (more liquidity/trades), then binance.us
	// Use combined stream for more frequent updates: miniTicker (1s) + aggTrade (every trade)
	endpoints := []string{
		"stream.binance.com:9443",
		"stream.binance.us:9443",
	}

	var lastErr error
	for _, host := range endpoints {
		wsURL := url.URL{
			Scheme:   "wss",
			Host:     host,
			Path:     "/stream",
			RawQuery: fmt.Sprintf("streams=%s@miniTicker/%s@aggTrade", c.symbol, c.symbol),
		}

		conn, _, err := dialer.DialContext(ctx, wsURL.String(), header)
		if err == nil {
			c.conn = conn
			log.Printf("connected to binance via %s", host)
			go c.readLoop(ctx)
			return nil
		}
		lastErr = err
		log.Printf("failed to connect to %s: %v, trying next...", host, err)
	}

	return fmt.Errorf("connect to binance (tried all endpoints): %w", lastErr)
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

	// Set up ping handler to keep connection alive
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start ping ticker to keep connection alive
	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			case <-pingTicker.C:
				c.mu.Lock()
				if c.conn != nil {
					c.conn.WriteMessage(websocket.PingMessage, nil)
				}
				c.mu.Unlock()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			log.Printf("binance read error: %v, reconnecting...", err)
			c.reconnect(ctx)
			continue
		}

		update, err := parseCombinedStream(message)
		if err != nil {
			log.Printf("parse stream error: %v", err)
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

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	header := make(map[string][]string)
	header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}

	endpoints := []string{
		"stream.binance.com:9443",
		"stream.binance.us:9443",
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		var lastErr error
		for _, host := range endpoints {
			wsURL := url.URL{
				Scheme:   "wss",
				Host:     host,
				Path:     "/stream",
				RawQuery: fmt.Sprintf("streams=%s@miniTicker/%s@aggTrade", c.symbol, c.symbol),
			}

			var err error
			c.conn, _, err = dialer.DialContext(ctx, wsURL.String(), header)
			if err == nil {
				log.Printf("binance reconnected via %s", host)
				return
			}
			lastErr = err
		}

		log.Printf("reconnect failed: %v, retrying in %v", lastErr, backoff)
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

// combinedStreamWrapper wraps messages from combined stream endpoint
type combinedStreamWrapper struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// miniTickerData represents 24hr mini ticker event data
type miniTickerData struct {
	EventType   string `json:"e"`
	EventTime   int64  `json:"E"`
	Symbol      string `json:"s"`
	ClosePrice  string `json:"c"`
	OpenPrice   string `json:"o"`
	HighPrice   string `json:"h"`
	LowPrice    string `json:"l"`
	BaseVolume  string `json:"v"`
	QuoteVolume string `json:"q"`
}

// aggTradeData represents aggregated trade event data
type aggTradeData struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int64  `json:"T"`
}

func parseCombinedStream(data []byte) (PriceUpdate, error) {
	var wrapper combinedStreamWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return PriceUpdate{}, fmt.Errorf("unmarshal wrapper: %w", err)
	}

	// Determine stream type from stream name
	if strings.Contains(wrapper.Stream, "@miniTicker") {
		var msg miniTickerData
		if err := json.Unmarshal(wrapper.Data, &msg); err != nil {
			return PriceUpdate{}, fmt.Errorf("unmarshal miniTicker: %w", err)
		}
		price, _ := strconv.ParseFloat(msg.ClosePrice, 64)
		volume, _ := strconv.ParseFloat(msg.BaseVolume, 64)
		timestamp := time.UnixMilli(msg.EventTime)
		if msg.EventTime == 0 {
			timestamp = time.Now()
		}
		return PriceUpdate{
			Symbol:    strings.ToUpper(msg.Symbol),
			Price:     price,
			Volume:    volume,
			Timestamp: timestamp,
		}, nil
	}

	if strings.Contains(wrapper.Stream, "@aggTrade") {
		var msg aggTradeData
		if err := json.Unmarshal(wrapper.Data, &msg); err != nil {
			return PriceUpdate{}, fmt.Errorf("unmarshal aggTrade: %w", err)
		}
		price, _ := strconv.ParseFloat(msg.Price, 64)
		quantity, _ := strconv.ParseFloat(msg.Quantity, 64)
		timestamp := time.UnixMilli(msg.TradeTime)
		if msg.TradeTime == 0 {
			timestamp = time.Now()
		}
		return PriceUpdate{
			Symbol:    strings.ToUpper(msg.Symbol),
			Price:     price,
			Volume:    quantity,
			Timestamp: timestamp,
		}, nil
	}

	return PriceUpdate{}, fmt.Errorf("unknown stream type: %s", wrapper.Stream)
}

// FetchKlines fetches historical candlestick data from Binance REST API.
// interval: 1m, 5m, 15m, 30m, 1h, 4h, 1d, etc.
// limit: number of candles to fetch (max 1000)
func FetchKlines(ctx context.Context, symbol, interval string, limit int) ([]Kline, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}

	// Try binance.us first, then binance.com
	endpoints := []string{
		"https://api.binance.us/api/v3/klines",
		"https://api.binance.com/api/v3/klines",
	}

	var lastErr error
	for _, baseURL := range endpoints {
		klines, err := fetchKlinesFromEndpoint(ctx, baseURL, symbol, interval, limit)
		if err == nil {
			return klines, nil
		}
		lastErr = err
		log.Printf("klines fetch from %s failed: %v, trying next", baseURL, err)
	}

	return nil, fmt.Errorf("all kline endpoints failed: %v", lastErr)
}

func fetchKlinesFromEndpoint(ctx context.Context, baseURL, symbol, interval string, limit int) ([]Kline, error) {
	reqURL := fmt.Sprintf("%s?symbol=%s&interval=%s&limit=%d",
		baseURL, strings.ToUpper(symbol), interval, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var rawKlines [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
		return nil, fmt.Errorf("decode klines: %w", err)
	}

	klines := make([]Kline, 0, len(rawKlines))
	for _, raw := range rawKlines {
		if len(raw) < 7 {
			continue
		}

		openTime, _ := raw[0].(float64)
		openStr, _ := raw[1].(string)
		highStr, _ := raw[2].(string)
		lowStr, _ := raw[3].(string)
		closeStr, _ := raw[4].(string)
		volumeStr, _ := raw[5].(string)
		closeTime, _ := raw[6].(float64)

		open, _ := strconv.ParseFloat(openStr, 64)
		high, _ := strconv.ParseFloat(highStr, 64)
		low, _ := strconv.ParseFloat(lowStr, 64)
		close, _ := strconv.ParseFloat(closeStr, 64)
		volume, _ := strconv.ParseFloat(volumeStr, 64)

		klines = append(klines, Kline{
			OpenTime:  time.UnixMilli(int64(openTime)),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
			CloseTime: time.UnixMilli(int64(closeTime)),
		})
	}

	log.Printf("fetched %d klines for %s from %s", len(klines), symbol, baseURL)
	return klines, nil
}
