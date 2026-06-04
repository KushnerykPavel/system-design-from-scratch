package main

import "testing"

func TestValidateAbusePlan(t *testing.T) {
	if issues := ValidateAbusePlan(AbusePlan{
		HasEdgeLimits:       true,
		HasAccountAwareGate: true,
		HasTenantBudget:     true,
		HasChallengeStep:    true,
	}); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	if issues := ValidateAbusePlan(AbusePlan{IPOnlyEnforcement: true}); len(issues) < 4 {
		t.Fatalf("expected several issues, got %v", issues)
	}
}
