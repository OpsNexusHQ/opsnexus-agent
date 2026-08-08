package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/OpsNexusHQ/opsnexus-agent/internal/collector"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/config"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/transport"
)

// Agent orchestrates config, collectors, and transport.
type Agent struct {
	cfg        *config.Config
	logger     *slog.Logger
	collectors []collector.Collector
	transport  transport.Transport
}

// NewAgent creates and configures a new Agent instance.
func NewAgent(cfg *config.Config, logger *slog.Logger) (*Agent, error) {
	// Initialize default collectors (e.g., system)
	sysCollector := collector.NewSystemCollector()

	// Initial console/logger transport placeholder
	consoleTransport := transport.NewConsoleTransport(logger)

	return &Agent{
		cfg:        cfg,
		logger:     logger,
		collectors: []collector.Collector{sysCollector},
		transport:  consoleTransport,
	}, nil
}

// Start runs the periodic metric collection and reporting loop.
// It exits when the provided context is canceled.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.Info("starting OpsNexus Agent", slog.String("interval", a.cfg.CollectionInterval))

	interval := a.cfg.ParseInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Trigger immediately at start
	a.collectAndSend(ctx)

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("agent collection loop stopped")
			return ctx.Err()
		case <-ticker.C:
			a.collectAndSend(ctx)
		}
	}
}

func (a *Agent) collectAndSend(ctx context.Context) {
	payload := make(map[string]any)

	for _, col := range a.collectors {
		select {
		case <-ctx.Done():
			return
		default:
		}

		data, err := col.Collect(ctx)
		if err != nil {
			a.logger.Error("failed to collect metrics", slog.String("collector", col.Name()), slog.Any("error", err))
			continue
		}
		payload[col.Name()] = data
	}

	if len(payload) > 0 {
		if err := a.transport.Send(ctx, payload); err != nil {
			a.logger.Error("failed to transmit metrics payload", slog.Any("error", err))
		}
	}
}
