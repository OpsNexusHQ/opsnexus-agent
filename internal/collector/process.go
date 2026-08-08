package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessMetrics holds the running process count.
type ProcessMetrics struct {
	RunningCount int `json:"running_count"`
}

// collectProcesses counts the number of running process IDs.
func collectProcesses(ctx context.Context) (*ProcessMetrics, error) {
	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("process pids: %w", err)
	}

	return &ProcessMetrics{RunningCount: len(pids)}, nil
}
