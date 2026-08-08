package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	envLogLevel           = "OPSNEXUS_AGENT_LOG_LEVEL"
	envCollectionInterval = "OPSNEXUS_AGENT_COLLECTION_INTERVAL"
	envBackendURL         = "OPSNEXUS_AGENT_BACKEND_URL"
)

// Config holds the configuration options for the OpsNexus Agent.
type Config struct {
	LogLevel           string `json:"log_level"`
	CollectionInterval string `json:"collection_interval"`
	BackendURL         string `json:"backend_url"`
}

// DefaultConfig returns a Config instance populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		LogLevel:           "info",
		CollectionInterval: "10s",
		BackendURL:         "http://localhost:8080",
	}
}

// ParseInterval returns the CollectionInterval parsed as a time.Duration.
func (c *Config) ParseInterval() (time.Duration, error) {
	d, err := time.ParseDuration(c.CollectionInterval)
	if err != nil {
		return 0, fmt.Errorf("invalid collection interval %q: %w", c.CollectionInterval, err)
	}

	if d <= 0 {
		return 0, errors.New("collection interval must be greater than zero")
	}

	return d, nil
}

// Validate performs basic configuration validation.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.BackendURL) == "" {
		return errors.New("backend_url is required")
	}

	parsed, err := url.Parse(strings.TrimSpace(c.BackendURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("backend_url must be a valid URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("backend_url must use http or https")
	}

	if _, err := c.ParseInterval(); err != nil {
		return err
	}

	return nil
}

// applyEnvOverrides applies environment variables on top of configuration values.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv(envLogLevel); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv(envCollectionInterval); v != "" {
		cfg.CollectionInterval = v
	}
	if v := os.Getenv(envBackendURL); v != "" {
		cfg.BackendURL = v
	}
}

// Load reads a configuration file from the specified path and overrides
// the defaults with any available environment variables.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
