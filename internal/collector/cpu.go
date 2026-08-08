package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUMetrics holds the collected CPU utilization data.
type CPUMetrics struct {
	UsagePercent float64   `json:"usage_percent"`
	PerCPU       []float64 `json:"per_cpu"`
	Count        int       `json:"count"`
}

// collectCPU gathers CPU utilization and core count.
// An interval of 0 computes usage since the last call,
// which suits periodic collection.
func collectCPU(ctx context.Context) (*CPUMetrics, error) {
	overall, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("cpu overall percent: %w", err)
	}

	var usagePercent float64
	if len(overall) > 0 {
		usagePercent = overall[0]
	}

	perCPU, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return nil, fmt.Errorf("cpu per-core percent: %w", err)
	}

	count, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("cpu count: %w", err)
	}

	return &CPUMetrics{
		UsagePercent: usagePercent,
		PerCPU:       perCPU,
		Count:        count,
	}, nil
}
