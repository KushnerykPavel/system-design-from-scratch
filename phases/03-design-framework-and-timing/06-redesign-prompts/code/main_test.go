package main

import "testing"

func TestRedesignPressure(t *testing.T) {
	changes := []ConstraintChange{
		{Name: "10x_qps", Severity: 3},
		{Name: "multi_region", Severity: 3},
		{Name: "lower_latency", Severity: 2},
	}

	if got := RedesignPressure(changes); got != 8 {
		t.Fatalf("expected pressure 8, got %d", got)
	}
	if !RequiresTopologyChange(changes) {
		t.Fatal("expected topology change requirement")
	}
}
