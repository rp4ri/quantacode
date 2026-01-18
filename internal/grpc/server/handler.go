package server

import (
	"log"
	"sync"
	"time"

	"github.com/yourusername/quantacode/internal/domain/indicators"
	"github.com/yourusername/quantacode/internal/infra/binance"
	pb "github.com/yourusername/quantacode/proto"
)

// Handler implements the MarketDataService gRPC server.
type Handler struct {
	pb.UnimplementedMarketDataServiceServer
	binanceClient *binance.Client
	mu            sync.RWMutex
}

// NewHandler creates a new gRPC handler with Binance client.
func NewHandler(binanceClient *binance.Client) *Handler {
	return &Handler{
		binanceClient: binanceClient,
	}
}

// StreamPrices implements bidirectional streaming of prices and indicators.
func (h *Handler) StreamPrices(req *pb.StreamRequest, stream pb.MarketDataService_StreamPricesServer) error {
	ctx := stream.Context()

	// Default indicator periods
	rsiPeriod := int(req.GetIndicators().GetRsiPeriod())
	smaPeriod := int(req.GetIndicators().GetSmaPeriod())
	emaPeriod := int(req.GetIndicators().GetEmaPeriod())

	if rsiPeriod <= 0 {
		rsiPeriod = 14
	}
	if smaPeriod <= 0 {
		smaPeriod = 14
	}
	if emaPeriod <= 0 {
		emaPeriod = 14
	}

	agg, err := indicators.NewAggregator(rsiPeriod, smaPeriod, emaPeriod)
	if err != nil {
		return err
	}

	priceCh := h.binanceClient.Subscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-priceCh:
			if !ok {
				return nil
			}

			// Send price update
			priceMsg := &pb.MarketUpdate{
				Update: &pb.MarketUpdate_Price{
					Price: &pb.PriceUpdate{
						Symbol:    update.Symbol,
						Price:     update.Price,
						Volume:    update.Volume,
						Timestamp: update.Timestamp.UnixMilli(),
					},
				},
			}
			if err := stream.Send(priceMsg); err != nil {
				log.Printf("send price error: %v", err)
				return err
			}

			// Calculate and send indicators
			vals := agg.Update(update.Price)
			history := agg.History()
			indicatorMsg := &pb.MarketUpdate{
				Update: &pb.MarketUpdate_Indicators{
					Indicators: &pb.IndicatorUpdate{
						Rsi:        vals.RSI,
						Sma:        vals.SMA,
						Ema:        vals.EMA,
						Timestamp:  time.Now().UnixMilli(),
						RsiHistory: history.RSI,
						SmaHistory: history.SMA,
						EmaHistory: history.EMA,
					},
				},
			}
			if err := stream.Send(indicatorMsg); err != nil {
				log.Printf("send indicators error: %v", err)
				return err
			}
		}
	}
}
