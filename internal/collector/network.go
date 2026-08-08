package collector

import (
	"context"
	"fmt"

	psnet "github.com/shirou/gopsutil/v4/net"
)

// NetworkMetrics holds I/O counters for all network interfaces.
type NetworkMetrics struct {
	Interfaces []InterfaceMetrics `json:"interfaces"`
}

// InterfaceMetrics holds I/O counters for a single network interface.
type InterfaceMetrics struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	ErrorsIn    uint64 `json:"errors_in"`
	ErrorsOut   uint64 `json:"errors_out"`
}

// collectNetwork gathers per-interface I/O counters.
func collectNetwork(ctx context.Context) (*NetworkMetrics, error) {
	counters, err := psnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("network io counters: %w", err)
	}

	interfaces := make([]InterfaceMetrics, 0, len(counters))
	for _, c := range counters {
		interfaces = append(interfaces, InterfaceMetrics{
			Name:        c.Name,
			BytesSent:   c.BytesSent,
			BytesRecv:   c.BytesRecv,
			PacketsSent: c.PacketsSent,
			PacketsRecv: c.PacketsRecv,
			ErrorsIn:    c.Errin,
			ErrorsOut:   c.Errout,
		})
	}

	return &NetworkMetrics{Interfaces: interfaces}, nil
}
