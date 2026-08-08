package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LogLevel != "info" {
		t.Errorf("expected LogLevel 'info', got %s", cfg.LogLevel)
	}

	if cfg.CollectionInterval != "10s" {
		t.Errorf("expected CollectionInterval '10s', got %s", cfg.CollectionInterval)
	}

	if cfg.BackendURL != "http://localhost:8080" {
		t.Errorf("expected BackendURL 'http://localhost:8080', got %s", cfg.BackendURL)
	}

	d, err := cfg.ParseInterval()
	if err != nil {
		t.Fatalf("expected valid interval, got %v", err)
	}
	if d != 10*time.Second {
		t.Errorf("expected parsed interval 10s, got %v", d)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error loading empty path: %v", err)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected LogLevel 'info', got %s", cfg.LogLevel)
	}

	if cfg.BackendURL != "http://localhost:8080" {
		t.Errorf("expected BackendURL 'http://localhost:8080', got %s", cfg.BackendURL)
	}
}

func TestLoadValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	content := []byte(`{
		"log_level": "debug",
		"collection_interval": "5s",
		"backend_url": "http://localhost:9090"
	}`)

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write tmp config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got %s", cfg.LogLevel)
	}

	if cfg.CollectionInterval != "5s" {
		t.Errorf("expected CollectionInterval '5s', got %s", cfg.CollectionInterval)
	}

	if cfg.BackendURL != "http://localhost:9090" {
		t.Errorf("expected BackendURL 'http://localhost:9090', got %s", cfg.BackendURL)
	}

	d, err := cfg.ParseInterval()
	if err != nil {
		t.Fatalf("expected valid interval, got %v", err)
	}
	if d != 5*time.Second {
		t.Errorf("expected parsed interval 5s, got %v", d)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	content := []byte(`{invalid-json}`)

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write tmp config file: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("expected error loading invalid JSON, got nil")
	}
}

func TestLoadFromEnvironmentOverrides(t *testing.T) {
	t.Setenv(envLogLevel, "debug")
	t.Setenv(envCollectionInterval, "15s")
	t.Setenv(envBackendURL, "http://localhost:9999")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got %s", cfg.LogLevel)
	}

	if cfg.CollectionInterval != "15s" {
		t.Errorf("expected CollectionInterval '15s', got %s", cfg.CollectionInterval)
	}

	if cfg.BackendURL != "http://localhost:9999" {
		t.Errorf("expected BackendURL 'http://localhost:9999', got %s", cfg.BackendURL)
	}
}

func TestParseIntervalRejectsInvalidDuration(t *testing.T) {
	cfg := &Config{CollectionInterval: "invalid", BackendURL: "http://localhost:8080"}
	if _, err := cfg.ParseInterval(); err == nil {
		t.Fatal("expected parse interval to fail for invalid duration")
	}
}
