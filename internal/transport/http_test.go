package transport

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHTTPTransportSendSuccess(t *testing.T) {
	var mu sync.Mutex
	var receivedBody []byte
	var receivedContentType string
	var receivedPath string

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		receivedPath = r.URL.Path
		receivedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		receivedBody = body
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	transport := NewHTTPTransport(backend.URL, "test-agent", logger)

	payload := map[string]any{
		"agent_id":  "test-agent",
		"timestamp": time.Now().Format(time.RFC3339),
		"metrics":   map[string]any{"cpu": 15.4},
	}

	err := transport.Send(context.Background(), payload)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if receivedPath != "/api/v1/agents/test-agent/telemetry" {
		t.Errorf("unexpected path: %s", receivedPath)
	}

	if receivedContentType != "application/json" {
		t.Errorf("unexpected content-type: %s", receivedContentType)
	}

	var decoded map[string]any
	if err := json.Unmarshal(receivedBody, &decoded); err != nil {
		t.Fatalf("failed to decode received body: %v", err)
	}

	if decoded["agent_id"] != "test-agent" {
		t.Errorf("expected agent_id 'test-agent', got %v", decoded["agent_id"])
	}
}

func TestHTTPTransportSendBackendFailure(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	transport := NewHTTPTransport(backend.URL, "test-agent", logger)

	err := transport.Send(context.Background(), map[string]any{"test": true})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestHTTPTransportSendContextCancellation(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusCreated)
	}))
	defer backend.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	transport := NewHTTPTransport(backend.URL, "test-agent", logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := transport.Send(ctx, map[string]any{"test": true})
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestHTTPTransportSendInvalidURL(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	transport := NewHTTPTransport("http://127.0.0.1:1", "test-agent", logger)

	err := transport.Send(context.Background(), map[string]any{"test": true})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
