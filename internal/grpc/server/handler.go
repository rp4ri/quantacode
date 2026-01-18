package server

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/yourusername/quantacode/internal/domain/indicators"
	"github.com/yourusername/quantacode/internal/infra/binance"
	pb "github.com/yourusername/quantacode/proto"
)

// Handler implements the MarketDataService gRPC server.
type Handler struct {
	pb.UnimplementedMarketDataServiceServer
	defaultSymbol string
	mu            sync.RWMutex
}

// NewHandler creates a new gRPC handler.
func NewHandler(defaultSymbol string) *Handler {
	return &Handler{
		defaultSymbol: strings.ToLower(defaultSymbol),
	}
}

// StreamPrices implements bidirectional streaming of prices and indicators.
func (h *Handler) StreamPrices(req *pb.StreamRequest, stream pb.MarketDataService_StreamPricesServer) error {
	ctx := stream.Context()

	// Use requested symbol or default
	symbol := strings.ToLower(req.GetSymbol())
	if symbol == "" {
		symbol = h.defaultSymbol
	}

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

	// CRITICAL: Fetch historical klines FIRST to pre-populate indicators
	// This ensures RSI/SMA/EMA are available from second 0
	klineCount := rsiPeriod + 10 // Fetch extra candles for accurate calculation
	if klineCount < 50 {
		klineCount = 50
	}
	klines, err := binance.FetchKlines(ctx, symbol, "1h", klineCount)
	if err != nil {
		log.Printf("warning: failed to fetch historical klines for %s: %v", symbol, err)
		// Continue anyway - indicators will warm up from real-time data
	} else {
		// Pre-populate aggregator with historical close prices
		for _, k := range klines {
			agg.Update(k.Close)
		}
		log.Printf("pre-populated indicators with %d historical candles for %s", len(klines), symbol)
	}

	// Create Binance client for requested symbol
	binanceClient := binance.NewClient(symbol)
	if err := binanceClient.Connect(ctx); err != nil {
		log.Printf("failed to connect to binance for %s: %v", symbol, err)
		return err
	}
	defer binanceClient.Close()

	// Send initial indicator values immediately (from historical data)
	if len(klines) > 0 {
		vals := agg.Values()
		history := agg.History()
		lastKline := klines[len(klines)-1]
		
		// Send initial price
		priceMsg := &pb.MarketUpdate{
			Update: &pb.MarketUpdate_Price{
				Price: &pb.PriceUpdate{
					Symbol:    strings.ToUpper(symbol),
					Price:     lastKline.Close,
					Volume:    lastKline.Volume,
					Timestamp: lastKline.CloseTime.UnixMilli(),
				},
			},
		}
		if err := stream.Send(priceMsg); err != nil {
			return err
		}
		
		// Send initial indicators
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
			return err
		}
		log.Printf("sent initial indicators for %s: RSI=%.2f SMA=%.2f EMA=%.2f", symbol, vals.RSI, vals.SMA, vals.EMA)
	}

	log.Printf("streaming %s for client", symbol)
	priceCh := binanceClient.Subscribe()

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
