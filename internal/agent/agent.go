package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/OpsNexusHQ/opsnexus-agent/internal/collector"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/config"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/transport"
	"github.com/OpsNexusHQ/opsnexus-common/models"
)

const agentID = "local-agent"

// Agent orchestrates configuration, collection, registration, and transport.
type Agent struct {
	cfg        *config.Config
	logger     *slog.Logger
	collectors []collector.Collector
	transport  transport.Transport
}

// NewAgent creates and configures a new Agent instance.
func NewAgent(cfg *config.Config, logger *slog.Logger) (*Agent, error) {
	sysCollector := collector.NewSystemCollector()

	httpTransport := transport.NewHTTPTransport(cfg.BackendURL, agentID, logger)

	return &Agent{
		cfg:        cfg,
		logger:     logger,
		collectors: []collector.Collector{sysCollector},
		transport:  httpTransport,
	}, nil
}

// Start registers the Agent and then starts the periodic
// metric collection and reporting loop.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.Info(
		"starting OpsNexus Agent",
		slog.String("interval", a.cfg.CollectionInterval),
		slog.String("backend_url", a.cfg.BackendURL),
	)

	if err := a.register(ctx); err != nil {
		return err
	}

	interval, err := a.cfg.ParseInterval()
	if err != nil {
		a.logger.Warn("invalid collection interval, falling back to default 10s",
			slog.String("collection_interval", a.cfg.CollectionInterval),
			slog.Any("error", err),
		)
		interval = 10 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Trigger immediately at start.
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

// register registers this Agent with the OpsNexus backend.
func (a *Agent) register(ctx context.Context) error {
	agent := models.Agent{
		ID:        agentID,
		Name:      "OpsNexus Local Agent",
		Hostname:  "localhost",
		OS:        "linux",
		Arch:      "amd64",
		Version:   "0.1.0",
		Status:    "online",
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
	}

	client := transport.NewRegistrationClient(a.cfg.BackendURL)

	if err := client.Register(ctx, agent); err != nil {
		a.logger.Error(
			"agent registration failed",
			slog.Any("error", err),
		)
		return err
	}

	a.logger.Info(
		"agent registered successfully",
		slog.String("agent_id", agent.ID),
	)

	return nil
}

func (a *Agent) collectAndSend(ctx context.Context) {
	metrics := make(map[string]any)

	for _, col := range a.collectors {
		select {
		case <-ctx.Done():
			return
		default:
		}

		data, err := col.Collect(ctx)
		if err != nil {
			a.logger.Error(
				"failed to collect metrics",
				slog.String("collector", col.Name()),
				slog.Any("error", err),
			)
			continue
		}

		metrics[col.Name()] = data
	}

	if len(metrics) == 0 {
		return
	}

	telemetry := models.AgentTelemetry{
		AgentID:   agentID,
		Timestamp: time.Now(),
		Metrics:   metrics,
	}

	if err := a.transport.Send(ctx, telemetry); err != nil {
		a.logger.Error(
			"telemetry transmission failed",
			slog.String("agent_id", agentID),
			slog.Any("error", err),
		)
	}
}
