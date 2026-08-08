package collector

import (
	"context"
	"testing"
)

func TestCollectMemory(t *testing.T) {
	metrics, err := collectMemory(context.Background())
	if err != nil {
		t.Fatalf("collectMemory failed: %v", err)
	}

	if metrics.TotalBytes == 0 {
		t.Error("expected non-zero TotalBytes")
	}

	if metrics.AvailableBytes == 0 {
		t.Error("expected non-zero AvailableBytes")
	}

	if metrics.UsedBytes == 0 {
		t.Error("expected non-zero UsedBytes")
	}

	if metrics.UsedPercent <= 0 || metrics.UsedPercent > 100 {
		t.Errorf("UsedPercent out of range: %f", metrics.UsedPercent)
	}
}
