package collector

import "context"

// Collector defines the interface that all metrics collectors must implement.
type Collector interface {
	Name() string
	Collect(ctx context.Context) (map[string]any, error)
}
