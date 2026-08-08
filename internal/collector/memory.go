package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

// MemoryMetrics holds virtual and swap memory statistics.
type MemoryMetrics struct {
	TotalBytes      uint64  `json:"total_bytes"`
	AvailableBytes  uint64  `json:"available_bytes"`
	UsedBytes       uint64  `json:"used_bytes"`
	UsedPercent     float64 `json:"used_percent"`
	SwapTotalBytes  uint64  `json:"swap_total_bytes"`
	SwapUsedBytes   uint64  `json:"swap_used_bytes"`
	SwapUsedPercent float64 `json:"swap_used_percent"`
}

// collectMemory gathers virtual memory and swap statistics.
func collectMemory(ctx context.Context) (*MemoryMetrics, error) {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("virtual memory: %w", err)
	}

	sm, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("swap memory: %w", err)
	}

	return &MemoryMetrics{
		TotalBytes:      vm.Total,
		AvailableBytes:  vm.Available,
		UsedBytes:       vm.Used,
		UsedPercent:     vm.UsedPercent,
		SwapTotalBytes:  sm.Total,
		SwapUsedBytes:   sm.Used,
		SwapUsedPercent: sm.UsedPercent,
	}, nil
}
