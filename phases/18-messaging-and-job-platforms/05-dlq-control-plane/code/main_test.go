package main

import "testing"

func TestValidateReplayRequest(t *testing.T) {
	req := ReplayRequest{
		HasReason:        true,
		HasScopedTarget:  true,
		HasTimeWindow:    true,
		HasDryRun:        true,
		HasRateLimit:     true,
		HasActorIdentity: true,
		HasRollbackPlan:  true,
	}
	if issues := ValidateReplayRequest(req); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	req.HasRateLimit = false
	req.HasDryRun = false
	if issues := ValidateReplayRequest(req); len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
}
