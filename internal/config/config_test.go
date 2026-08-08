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
	if cfg.ParseInterval() != 10*time.Second {
		t.Errorf("expected parsed interval 10s, got %v", cfg.ParseInterval())
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
}

func TestLoadValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	content := []byte(`{"log_level": "debug", "collection_interval": "5s"}`)
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
	if cfg.ParseInterval() != 5*time.Second {
		t.Errorf("expected parsed interval 5s, got %v", cfg.ParseInterval())
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
