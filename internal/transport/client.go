package transport

import (
	"context"
	"log/slog"
)

// ConsoleTransport prints payloads using the structured logger.
// It serves as the initial foundation placeholder for transport.
type ConsoleTransport struct {
	logger *slog.Logger
}

// NewConsoleTransport creates a new ConsoleTransport instance.
func NewConsoleTransport(logger *slog.Logger) *ConsoleTransport {
	return &ConsoleTransport{
		logger: logger,
	}
}

// Send logs the payload as structured data, satisfying the Transport interface.
func (t *ConsoleTransport) Send(ctx context.Context, payload any) error {
	t.logger.Info("transport payload transmitted", slog.Any("payload", payload))
	return nil
}
var _ Transport = (*ConsoleTransport)(nil)
