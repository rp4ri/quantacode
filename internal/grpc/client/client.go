package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/yourusername/quantacode/proto"
)

// PriceUpdate represents a price tick.
type PriceUpdate struct {
	Symbol    string
	Price     float64
	Volume    float64
	Timestamp time.Time
}

// IndicatorUpdate represents indicator values.
type IndicatorUpdate struct {
	RSI        float64
	SMA        float64
	EMA        float64
	Timestamp  time.Time
	RSIHistory []float64
	SMAHistory []float64
	EMAHistory []float64
}

// Client manages gRPC connection to the server.
type Client struct {
	conn   *grpc.ClientConn
	client pb.MarketDataServiceClient
}

// New creates a new gRPC client.
func New(serverAddr string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to server: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewMarketDataServiceClient(conn),
	}, nil
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// StreamPrices starts streaming prices and indicators.
func (c *Client) StreamPrices(ctx context.Context, symbol string, priceCh chan<- PriceUpdate, indicatorCh chan<- IndicatorUpdate) error {
	req := &pb.StreamRequest{
		Symbol: symbol,
		Indicators: &pb.IndicatorConfig{
			RsiPeriod: 14,
			SmaPeriod: 14,
			EmaPeriod: 14,
		},
	}

	stream, err := c.client.StreamPrices(ctx, req)
	if err != nil {
		return fmt.Errorf("start stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("receive: %w", err)
		}

		switch update := msg.Update.(type) {
		case *pb.MarketUpdate_Price:
			priceCh <- PriceUpdate{
				Symbol:    update.Price.Symbol,
				Price:     update.Price.Price,
				Volume:    update.Price.Volume,
				Timestamp: time.UnixMilli(update.Price.Timestamp),
			}
		case *pb.MarketUpdate_Indicators:
			indicatorCh <- IndicatorUpdate{
				RSI:        update.Indicators.Rsi,
				SMA:        update.Indicators.Sma,
				EMA:        update.Indicators.Ema,
				Timestamp:  time.UnixMilli(update.Indicators.Timestamp),
				RSIHistory: update.Indicators.RsiHistory,
				SMAHistory: update.Indicators.SmaHistory,
				EMAHistory: update.Indicators.EmaHistory,
			}
		default:
			log.Printf("unknown update type: %T", update)
		}
	}
}
