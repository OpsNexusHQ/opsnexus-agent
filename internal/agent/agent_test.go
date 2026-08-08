package agent

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/OpsNexusHQ/opsnexus-agent/internal/config"
)

func TestAgentLifecycle(t *testing.T) {
	// Discard log output in tests
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		LogLevel:           "info",
		CollectionInterval: "50ms",
	}

	a, err := NewAgent(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err = a.Start(ctx)
	if err != context.DeadlineExceeded && err != nil {
		t.Errorf("expected context deadline exceeded or nil, got: %v", err)
	}
}
