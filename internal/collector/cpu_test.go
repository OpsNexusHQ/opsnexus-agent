package collector

import (
	"context"
	"testing"
)

func TestCollectCPU(t *testing.T) {
	metrics, err := collectCPU(context.Background())
	if err != nil {
		t.Fatalf("collectCPU failed: %v", err)
	}

	if metrics.Count <= 0 {
		t.Errorf("expected positive CPU count, got %d", metrics.Count)
	}

	if len(metrics.PerCPU) == 0 {
		t.Error("expected non-empty PerCPU slice")
	}

	if len(metrics.PerCPU) != metrics.Count {
		t.Errorf("PerCPU length (%d) does not match Count (%d)", len(metrics.PerCPU), metrics.Count)
	}
}
