package main

import "testing"

func TestValidateWorkflowPlan(t *testing.T) {
	plan := WorkflowPlan{
		HasDurableState:         true,
		HasTimerQueue:           true,
		HasIdempotentActivities: true,
		SupportsCompensation:    true,
		HasHeartbeats:           true,
		HasVersioning:           true,
		HasIsolationControls:    true,
	}
	if issues := ValidateWorkflowPlan(plan); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	plan.HasVersioning = false
	plan.HasTimerQueue = false
	if issues := ValidateWorkflowPlan(plan); len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
}
