package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// HTTPTransport sends telemetry payloads to the OpsNexus backend via HTTP.
type HTTPTransport struct {
	backendURL string
	agentID    string
	logger     *slog.Logger
	httpClient *http.Client
}

// NewHTTPTransport creates an HTTP transport targeting the backend telemetry endpoint.
func NewHTTPTransport(backendURL, agentID string, logger *slog.Logger) *HTTPTransport {
	return &HTTPTransport{
		backendURL: strings.TrimRight(backendURL, "/"),
		agentID:    agentID,
		logger:     logger,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send serializes the payload to JSON and POSTs it to the backend
// telemetry ingestion endpoint. Errors are returned but never cause panics.
func (t *HTTPTransport) Send(ctx context.Context, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telemetry payload: %w", err)
	}

	url := t.backendURL + "/api/v1/agents/" + t.agentID + "/telemetry"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create telemetry request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send telemetry request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telemetry rejected with HTTP status %d", resp.StatusCode)
	}

	t.logger.Debug("telemetry transmitted successfully",
		slog.String("agent_id", t.agentID),
		slog.Int("status", resp.StatusCode),
	)

	return nil
}

var _ Transport = (*HTTPTransport)(nil)
