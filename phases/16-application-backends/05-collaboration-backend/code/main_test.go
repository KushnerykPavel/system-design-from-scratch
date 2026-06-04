package main

import "testing"

func TestValidateCollaborationPlanHealthy(t *testing.T) {
	plan := CollaborationPlan{
		Name:                  "healthy",
		ConvergenceMode:       "sequenced",
		HasSessionOwner:       true,
		HasSnapshots:          true,
		SnapshotIntervalOps:   500,
		HasPresenceTTL:        true,
		SupportsReplay:        true,
		HasDeterministicCheck: true,
	}
	if issues := ValidateCollaborationPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateCollaborationPlan returned issues: %v", issues)
	}
}

func TestValidateCollaborationPlanRejectsWeakPlan(t *testing.T) {
	plan := CollaborationPlan{
		Name:            "weak",
		ConvergenceMode: "unknown",
	}
	if issues := ValidateCollaborationPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateCollaborationPlan returned too few issues: %v", issues)
	}
}
