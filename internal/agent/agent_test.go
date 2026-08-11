package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/OpsNexusHQ/opsnexus-agent/internal/collector"
	"github.com/OpsNexusHQ/opsnexus-agent/internal/config"
	"github.com/OpsNexusHQ/opsnexus-common/models"
)

type mockCollector struct {
	mu    sync.Mutex
	calls int
}

func (m *mockCollector) Name() string {
	return "mock"
}

func (m *mockCollector) Collect(ctx context.Context) (map[string]any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls++
	return map[string]any{"value": 42}, nil
}

func (m *mockCollector) getCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

type mockTransport struct {
	mu       sync.Mutex
	payloads []any
}

func (m *mockTransport) Send(ctx context.Context, payload any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.payloads = append(m.payloads, payload)
	return nil
}

func (m *mockTransport) getPayloads() []any {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]any, len(m.payloads))
	copy(res, m.payloads)
	return res
}

func TestAgentLifecycle(t *testing.T) {
	// Create a fake backend for the registration request.
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/agents/register" {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		// Accept telemetry POSTs
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}))
	defer backend.Close()

	// Discard log output in tests.
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &config.Config{
		LogLevel:           "info",
		CollectionInterval: "50ms",
		BackendURL:         backend.URL,
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

func TestAgentCollectionIntervalAndCancellation(t *testing.T) {
	// Create a fake backend for registration
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		LogLevel:           "info",
		CollectionInterval: "20ms",
		BackendURL:         backend.URL,
	}

	a, err := NewAgent(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	// Inject mocks
	mc := &mockCollector{}
	mt := &mockTransport{}
	a.collectors = []collector.Collector{mc}
	a.transport = mt

	// Run agent with timeout to verify cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	err = a.Start(ctx)
	if err != context.DeadlineExceeded && err != nil {
		t.Errorf("expected context deadline exceeded or nil, got: %v", err)
	}

	calls := mc.getCalls()
	if calls < 3 {
		t.Errorf("expected at least 3 collector invocations, got %d", calls)
	}

	payloads := mt.getPayloads()
	if len(payloads) != calls {
		t.Errorf("expected payload count to match calls (%d), got %d", calls, len(payloads))
	}
}

func TestAgentConstructsAgentTelemetry(t *testing.T) {
	// Create a fake backend for registration
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		LogLevel:           "info",
		CollectionInterval: "20ms",
		BackendURL:         backend.URL,
	}

	a, err := NewAgent(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	// Inject mock collector and transport
	mc := &mockCollector{}
	mt := &mockTransport{}
	a.collectors = []collector.Collector{mc}
	a.transport = mt

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	_ = a.Start(ctx)

	payloads := mt.getPayloads()
	if len(payloads) == 0 {
		t.Fatal("expected at least one payload")
	}

	// Verify the payload is an AgentTelemetry struct
	telemetry, ok := payloads[0].(models.AgentTelemetry)
	if !ok {
		t.Fatalf("expected payload to be models.AgentTelemetry, got %T", payloads[0])
	}

	if telemetry.AgentID != "local-agent" {
		t.Errorf("expected agent_id 'local-agent', got %s", telemetry.AgentID)
	}

	if telemetry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	if telemetry.Metrics == nil {
		t.Fatal("expected non-nil metrics")
	}

	// Verify the mock collector data is nested under "mock"
	mockData, ok := telemetry.Metrics["mock"]
	if !ok {
		t.Fatal("expected key 'mock' in metrics")
	}

	dataMap, ok := mockData.(map[string]any)
	if !ok {
		t.Fatalf("expected mock data to be map[string]any, got %T", mockData)
	}

	if dataMap["value"] != 42 {
		t.Errorf("expected value 42, got %v", dataMap["value"])
	}
}

func TestAgentLogsSuccessfulCollectionAndTransmission(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	a := &Agent{
		logger:     logger,
		collectors: []collector.Collector{&mockCollector{}},
		transport:  &mockTransport{},
	}

	a.collectAndSend(context.Background())
	output := logs.String()
	for _, message := range []string{
		"msg=\"metrics collected successfully\"",
		"collector=mock",
		"duration=",
		"msg=\"telemetry transmitted successfully\"",
		"metric_sources=1",
	} {
		if !strings.Contains(output, message) {
			t.Errorf("expected log output to contain %q, got %q", message, output)
		}
	}
}

func TestAgentSendsFullSystemTelemetry(t *testing.T) {
	var gotPayload []byte
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/agents/register" {
			w.WriteHeader(http.StatusCreated)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected telemetry POST, got %s", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		gotPayload = body
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		LogLevel:           "info",
		CollectionInterval: "20ms",
		BackendURL:         backend.URL,
	}

	a, err := NewAgent(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	_ = a.Start(ctx)

	if len(gotPayload) == 0 {
		t.Fatal("expected telemetry payload to be sent")
	}

	var decoded map[string]any
	if err := json.Unmarshal(gotPayload, &decoded); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	metrics, ok := decoded["metrics"].(map[string]any)
	if !ok {
		t.Fatal("expected metrics to be an object")
	}

	system, ok := metrics["system"].(map[string]any)
	if !ok {
		t.Fatal("expected system metrics to be an object")
	}

	if _, ok := system["network"]; !ok {
		t.Fatal("expected network data in system metrics")
	}

	if _, ok := system["processes"]; !ok {
		t.Fatal("expected processes data in system metrics")
	}
}
