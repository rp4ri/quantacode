package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/yourusername/quantacode/internal/grpc/server"
	"github.com/yourusername/quantacode/internal/infra/binance"
	pb "github.com/yourusername/quantacode/proto"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		symbol = "btcusdt"
	}
	symbol = strings.ToLower(symbol)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize Binance client (try real, fallback to simulated)
	binanceClient := binance.NewClient(symbol)
	if err := binanceClient.Connect(ctx); err != nil {
		log.Printf("failed to connect to binance: %v, using simulated data", err)
		binanceClient = binance.NewSimulatedClient(symbol)
		if err := binanceClient.Connect(ctx); err != nil {
			log.Fatalf("failed to start simulated client: %v", err)
		}
	}
	defer binanceClient.Close()

	// Create gRPC server
	handler := server.NewHandler(binanceClient)
	grpcServer := grpc.NewServer()
	pb.RegisterMarketDataServiceServer(grpcServer, handler)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Graceful shutdown goroutine
	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")

		// Give 5 seconds for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("server stopped gracefully")
		case <-shutdownCtx.Done():
			log.Println("forcing server stop")
			grpcServer.Stop()
		}
	}()

	log.Printf("server listening on :%s (symbol: %s)", port, symbol)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("server error: %v", err)
	}
}
