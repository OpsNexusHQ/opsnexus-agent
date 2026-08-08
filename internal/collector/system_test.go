package collector

import (
	"context"
	"testing"
)

func TestSystemCollectorInterface(t *testing.T) {
	var _ Collector = (*SystemCollector)(nil)
}

func TestSystemCollectorName(t *testing.T) {
	c := NewSystemCollector()
	if c.Name() != "system" {
		t.Errorf("expected collector name 'system', got %s", c.Name())
	}
}

func TestSystemCollectorCollect(t *testing.T) {
	c := NewSystemCollector()

	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error during collection: %v", err)
	}

	if metrics == nil {
		t.Fatal("expected non-nil metrics map")
	}

	expectedKeys := []string{"cpu", "memory", "disk", "network", "uptime", "processes", "timestamp"}
	for _, key := range expectedKeys {
		if _, ok := metrics[key]; !ok {
			t.Errorf("expected key %q in metrics map", key)
		}
	}
}

func TestSystemSnapshotCPU(t *testing.T) {
	c := NewSystemCollector()
	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cpuData, ok := metrics["cpu"]
	if !ok {
		t.Fatal("missing 'cpu' in metrics")
	}

	cpuMetrics, ok := cpuData.(*CPUMetrics)
	if !ok {
		t.Fatalf("cpu metrics unexpected type: %T", cpuData)
	}

	if cpuMetrics.Count <= 0 {
		t.Errorf("expected positive CPU count, got %d", cpuMetrics.Count)
	}
}

func TestSystemSnapshotMemory(t *testing.T) {
	c := NewSystemCollector()
	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	memData, ok := metrics["memory"]
	if !ok {
		t.Fatal("missing 'memory' in metrics")
	}

	memMetrics, ok := memData.(*MemoryMetrics)
	if !ok {
		t.Fatalf("memory metrics unexpected type: %T", memData)
	}

	if memMetrics.TotalBytes == 0 {
		t.Error("expected non-zero TotalBytes")
	}
}

func TestSystemSnapshotUptime(t *testing.T) {
	c := NewSystemCollector()
	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	uptimeData, ok := metrics["uptime"]
	if !ok {
		t.Fatal("missing 'uptime' in metrics")
	}

	uptimeMetrics, ok := uptimeData.(*UptimeMetrics)
	if !ok {
		t.Fatalf("uptime metrics unexpected type: %T", uptimeData)
	}

	if uptimeMetrics.UptimeSeconds == 0 {
		t.Error("expected non-zero uptime")
	}
}

func TestSystemSnapshotProcesses(t *testing.T) {
	c := NewSystemCollector()
	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	procData, ok := metrics["processes"]
	if !ok {
		t.Fatal("missing 'processes' in metrics")
	}

	procMetrics, ok := procData.(*ProcessMetrics)
	if !ok {
		t.Fatalf("process metrics unexpected type: %T", procData)
	}

	if procMetrics.RunningCount <= 0 {
		t.Errorf("expected positive process count, got %d", procMetrics.RunningCount)
	}
}
