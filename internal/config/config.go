package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the configuration options for the OpsNexus Agent.
type Config struct {
	LogLevel           string `json:"log_level"`
	CollectionInterval string `json:"collection_interval"`
}

// DefaultConfig returns a Config instance populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		LogLevel:           "info",
		CollectionInterval: "10s",
	}
}

// ParseInterval returns the CollectionInterval parsed as a time.Duration.
// If the interval is invalid, it returns a default duration of 10 seconds.
func (c *Config) ParseInterval() time.Duration {
	d, err := time.ParseDuration(c.CollectionInterval)
	if err != nil {
		return 10 * time.Second
	}
	return d
}

// Load reads a configuration file from the specified path and overrides
// the defaults. If path is empty, it returns the defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
