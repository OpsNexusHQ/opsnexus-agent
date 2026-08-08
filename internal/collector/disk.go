package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

// DiskMetrics holds usage data for all discovered partitions.
type DiskMetrics struct {
	Partitions []PartitionMetrics `json:"partitions"`
}

// PartitionMetrics holds usage data for a single disk partition.
type PartitionMetrics struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

// collectDisk gathers partition and usage information.
// Partitions that cannot be read (e.g. permission denied) are silently skipped.
func collectDisk(ctx context.Context) (*DiskMetrics, error) {
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("disk partitions: %w", err)
	}

	var metrics []PartitionMetrics
	for _, p := range partitions {
		usage, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil {
			// Skip partitions we can't read.
			continue
		}

		metrics = append(metrics, PartitionMetrics{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return &DiskMetrics{Partitions: metrics}, nil
}
