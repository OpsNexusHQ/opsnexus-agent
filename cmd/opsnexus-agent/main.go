package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/OpsNexusHQ/opsnexus-agent/internal/agent"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to the agent configuration file")
	flag.Parse()

	// 1. Load Configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", os.ErrNotExist)
		os.Exit(1)
	}

	// 2. Set Up Structured Logging
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	logger.Info("opsnexus-agent initializing")

	// 3. Setup Context with Graceful Shutdown Handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Initialize Agent
	a, err := agent.NewAgent(cfg, logger)
	if err != nil {
		logger.Error("failed to create agent", slog.Any("error", err))
		os.Exit(1)
	}

	// 5. Start Agent
	go func() {
		if err := a.Start(ctx); err != nil && err != context.Canceled {
			logger.Error("agent stopped with error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("shutting down opsnexus-agent gracefully")
}
