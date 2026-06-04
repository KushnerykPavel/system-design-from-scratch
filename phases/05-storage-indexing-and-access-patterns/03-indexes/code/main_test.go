package main

import "testing"

func TestWriteUnitsPerLogicalWrite(t *testing.T) {
	plan := IndexPlan{
		BaseWriteCost:     2,
		IndexCount:        3,
		ReplicationFactor: 2,
	}

	got := plan.WriteUnitsPerLogicalWrite()
	if got != 16 {
		t.Fatalf("expected 16 write units, got %d", got)
	}
}
