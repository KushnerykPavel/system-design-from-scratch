package main

import "testing"

func TestAssessThreatModel(t *testing.T) {
	strong := AssessThreatModel(ThreatModel{
		HasAssets:        true,
		HasActors:        true,
		HasBoundaries:    true,
		HasTopThreats:    true,
		HasMitigations:   true,
		HasObservability: true,
		HasDesignChange:  true,
	})
	if strong.Level != "strong" {
		t.Fatalf("expected strong, got %+v", strong)
	}

	weak := AssessThreatModel(ThreatModel{})
	if weak.Level == "strong" || len(weak.Missing) == 0 {
		t.Fatalf("expected weak with missing items, got %+v", weak)
	}
}
