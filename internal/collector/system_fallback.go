//go:build !linux

package collector

import (
	"context"
)

// Collect returns dummy/fallback metrics on non-Linux platforms.
func (s *SystemCollector) Collect(ctx context.Context) (map[string]any, error) {
	return map[string]any{
		"cpu": map[string]float64{
			"utilization_percent": 15.4,
		},
		"load": map[string]float64{
			"load1":  0.42,
			"load5":  0.55,
			"load15": 0.61,
		},
		"memory": map[string]uint64{
			"total_bytes":     16000000000,
			"free_bytes":      4000000000,
			"available_bytes": 10000000000,
		},
	}, nil
}
