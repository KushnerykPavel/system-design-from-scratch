package main

import "testing"

func TestValidateRingPlanAcceptsReasonablePlan(t *testing.T) {
	plan := RingPlan{
		Name:                    "good",
		Nodes:                   50,
		VirtualNodesPerNode:     64,
		ExpectedRemapPercent:    10,
		WeightedPlacement:       true,
		HealthAwareRouting:      true,
		PrewarmOnRebalance:      true,
		HotKeyIsolationStrategy: true,
	}
	if issues := ValidateRingPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateRingPlan() returned issues: %v", issues)
	}
}

func TestValidateRingPlanRejectsWeakPlan(t *testing.T) {
	plan := RingPlan{
		Name:                    "bad",
		Nodes:                   1,
		VirtualNodesPerNode:     4,
		ExpectedRemapPercent:    50,
		WeightedPlacement:       false,
		HealthAwareRouting:      false,
		PrewarmOnRebalance:      false,
		HotKeyIsolationStrategy: false,
	}
	if issues := ValidateRingPlan(plan); len(issues) < 4 {
		t.Fatalf("ValidateRingPlan() returned too few issues: %v", issues)
	}
}
