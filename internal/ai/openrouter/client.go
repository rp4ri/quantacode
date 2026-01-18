package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/quantacode/internal/logging"
)

const (
	baseURL      = "https://openrouter.ai/api/v1/chat/completions"
	defaultModel = "deepseek/deepseek-chat"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	model      string
	logger     *logging.Logger
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Delta        *Delta  `json:"delta,omitempty"`
	Message      *Delta  `json:"message,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StreamChunk struct {
	Content      string
	Done         bool
	Error        error
	FinishReason string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		model:  defaultModel,
		logger: logging.GetLogger("openrouter"),
	}
}

func (c *Client) WithModel(model string) *Client {
	c.model = model
	return c
}

func (c *Client) buildSystemPrompt(symbol string, price float64, rsi, sma, ema float64) string {
	return fmt.Sprintf(`Eres un analista de trading experto especializado en criptomonedas. 
Estás analizando %s en tiempo real.

Datos actuales del mercado:
- Precio: $%.2f
- RSI (14): %.2f
- SMA (14): %.2f  
- EMA (14): %.2f

Proporciona análisis técnico conciso y accionable. Responde en español.
Sé directo y específico sobre señales de compra/venta basándote en los indicadores.`, 
		symbol, price, rsi, sma, ema)
}

func (c *Client) StreamAnalysis(ctx context.Context, userPrompt, symbol string, price, rsi, sma, ema float64) (<-chan StreamChunk, error) {
	systemPrompt := c.buildSystemPrompt(symbol, price, rsi, sma, ema)
	
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Stream:      true,
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	c.logger.LogIO("OpenRouter request", map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
		"symbol":   symbol,
		"price":    price,
		"rsi":      rsi,
		"sma":      sma,
		"ema":      ema,
	}, nil, 0)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://quantacode.local")
	httpReq.Header.Set("X-Title", "QuantaCode Trading CLI")

	startTime := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.LogOpenRouterCall(req, nil, err, time.Since(startTime))
		return nil, fmt.Errorf("http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		err := fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
		c.logger.LogOpenRouterCall(req, string(body), err, time.Since(startTime))
		return nil, err
	}

	chunkCh := make(chan StreamChunk, 100)

	go func() {
		defer resp.Body.Close()
		defer close(chunkCh)

		var fullResponse strings.Builder
		reader := bufio.NewReader(resp.Body)

		for {
			select {
			case <-ctx.Done():
				chunkCh <- StreamChunk{Error: ctx.Err(), Done: true}
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					chunkCh <- StreamChunk{Done: true}
					c.logger.LogOpenRouterCall(req, fullResponse.String(), nil, time.Since(startTime))
					return
				}
				chunkCh <- StreamChunk{Error: err, Done: true}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				chunkCh <- StreamChunk{Done: true}
				c.logger.LogOpenRouterCall(req, fullResponse.String(), nil, time.Since(startTime))
				return
			}

			var chatResp ChatResponse
			if err := json.Unmarshal([]byte(data), &chatResp); err != nil {
				c.logger.Error("parse SSE chunk", err)
				continue
			}

			if chatResp.Error != nil {
				chunkCh <- StreamChunk{
					Error: fmt.Errorf("%s: %s", chatResp.Error.Code, chatResp.Error.Message),
					Done:  true,
				}
				return
			}

			for _, choice := range chatResp.Choices {
				if choice.Delta != nil && choice.Delta.Content != "" {
					fullResponse.WriteString(choice.Delta.Content)
					chunkCh <- StreamChunk{
						Content:      choice.Delta.Content,
						FinishReason: choice.FinishReason,
					}
				}
				if choice.FinishReason != "" && choice.FinishReason != "null" {
					chunkCh <- StreamChunk{Done: true, FinishReason: choice.FinishReason}
					c.logger.LogOpenRouterCall(req, fullResponse.String(), nil, time.Since(startTime))
					return
				}
			}
		}
	}()

	return chunkCh, nil
}
