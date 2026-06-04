package main

import "testing"

func TestValidateIsolationPlan(t *testing.T) {
	if issues := ValidateIsolationPlan(IsolationPlan{
		TenantScopedKeys: true,
		FairScheduling:   true,
		PerTenantBudgets: true,
		CarveOutPath:     true,
	}); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	if issues := ValidateIsolationPlan(IsolationPlan{SharedAdminBypass: true}); len(issues) < 4 {
		t.Fatalf("expected several issues, got %v", issues)
	}
}
