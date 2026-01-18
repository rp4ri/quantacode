package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/yourusername/quantacode/internal/ui/chat"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "quantacode",
		Short: "QuantaCode CLI for streaming market analysis",
	}

	root.AddCommand(newChatCmd())
	return root
}

func newChatCmd() *cobra.Command {
	var (
		serverAddr string
		symbol     string
		keyFlag    string
	)

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Open interactive market analysis chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyFlag == "" {
				keyFlag = os.Getenv("OPENROUTER_API_KEY")
			}
			if keyFlag == "" {
				keyFlag = os.Getenv("OPENROUTER_KEY")
			}
			if keyFlag == "" {
				return fmt.Errorf("OpenRouter API key not provided (use --openrouter-key or set OPENROUTER_API_KEY)")
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			cfg := chat.Config{
				ServerAddr:    serverAddr,
				Symbol:        symbol,
				OpenRouterKey: keyFlag,
			}
			return chat.Run(ctx, cfg)
		},
	}

	cmd.Flags().StringVar(&serverAddr, "server", "localhost:50051", "gRPC server address")
	cmd.Flags().StringVar(&symbol, "symbol", "BTCUSDT", "Trading symbol to subscribe to")
	cmd.Flags().StringVar(&keyFlag, "openrouter-key", "", "OpenRouter API key (fallback to OPENROUTER_KEY env var)")

	return cmd
}
