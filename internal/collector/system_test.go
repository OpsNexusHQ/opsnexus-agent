package collector

import (
	"context"
	"testing"
)

func TestSystemCollector(t *testing.T) {
	c := NewSystemCollector()
	if c.Name() != "system" {
		t.Errorf("expected collector name 'system', got %s", c.Name())
	}

	metrics, err := c.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error during collection: %v", err)
	}

	if metrics == nil {
		t.Fatal("expected non-nil metrics map")
	}

	// Verify we got the basic metric structures
	if _, ok := metrics["cpu"]; !ok {
		t.Error("expected 'cpu' metrics to be present")
	}
	if _, ok := metrics["load"]; !ok {
		t.Error("expected 'load' metrics to be present")
	}
	if _, ok := metrics["memory"]; !ok {
		t.Error("expected 'memory' metrics to be present")
	}
}
