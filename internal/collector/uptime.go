package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
)

// UptimeMetrics holds system uptime data.
type UptimeMetrics struct {
	UptimeSeconds uint64 `json:"uptime_seconds"`
}

// collectUptime gathers system uptime in seconds.
func collectUptime(ctx context.Context) (*UptimeMetrics, error) {
	uptime, err := host.UptimeWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("host uptime: %w", err)
	}

	return &UptimeMetrics{UptimeSeconds: uptime}, nil
}
