package collector

import (
	"context"
	"log/slog"
	"time"
)

// SystemSnapshot is the aggregated point-in-time system state
// returned by SystemCollector.Collect.
type SystemSnapshot struct {
	CPU       *CPUMetrics     `json:"cpu,omitempty"`
	Memory    *MemoryMetrics  `json:"memory,omitempty"`
	Disk      *DiskMetrics    `json:"disk,omitempty"`
	Network   *NetworkMetrics `json:"network,omitempty"`
	Uptime    *UptimeMetrics  `json:"uptime,omitempty"`
	Processes *ProcessMetrics `json:"processes,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// SystemCollector collects system metrics and assembles a SystemSnapshot.
type SystemCollector struct{}

// NewSystemCollector creates a new instance of SystemCollector.
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

// Name returns the identifier of this collector.
func (s *SystemCollector) Name() string {
	return "system"
}

// Collect gathers all system metrics into a SystemSnapshot
// and returns it as a map for Collector interface compatibility.
// Individual sub-collector failures are logged and skipped;
// partial results are returned rather than failing entirely.
func (s *SystemCollector) Collect(ctx context.Context) (map[string]any, error) {
	snapshot := &SystemSnapshot{
		Timestamp: time.Now(),
	}

	logger := slog.Default()

	if cpuMetrics, err := collectCPU(ctx); err != nil {
		logger.Warn("cpu collection failed", slog.Any("error", err))
	} else {
		snapshot.CPU = cpuMetrics
	}

	if memMetrics, err := collectMemory(ctx); err != nil {
		logger.Warn("memory collection failed", slog.Any("error", err))
	} else {
		snapshot.Memory = memMetrics
	}

	if diskMetrics, err := collectDisk(ctx); err != nil {
		logger.Warn("disk collection failed", slog.Any("error", err))
	} else {
		snapshot.Disk = diskMetrics
	}

	if netMetrics, err := collectNetwork(ctx); err != nil {
		logger.Warn("network collection failed", slog.Any("error", err))
	} else {
		snapshot.Network = netMetrics
	}

	if uptimeMetrics, err := collectUptime(ctx); err != nil {
		logger.Warn("uptime collection failed", slog.Any("error", err))
	} else {
		snapshot.Uptime = uptimeMetrics
	}

	if procMetrics, err := collectProcesses(ctx); err != nil {
		logger.Warn("process collection failed", slog.Any("error", err))
	} else {
		snapshot.Processes = procMetrics
	}

	return map[string]any{
		"cpu":       snapshot.CPU,
		"memory":    snapshot.Memory,
		"disk":      snapshot.Disk,
		"network":   snapshot.Network,
		"uptime":    snapshot.Uptime,
		"processes": snapshot.Processes,
		"timestamp": snapshot.Timestamp,
	}, nil
}

// Verify SystemCollector satisfies the Collector interface.
var _ Collector = (*SystemCollector)(nil)
