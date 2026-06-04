package main

import "testing"

func TestValidateLifecyclePolicyAcceptsReasonablePolicy(t *testing.T) {
	cfg := LifecyclePolicy{
		Name:                "reasonable",
		TransitionAfterDays: 30,
		ExpireAfterDays:     365,
		DeleteGraceDays:     7,
		LegalHoldAware:      true,
		DryRunRequired:      true,
		CompactionBudgetPct: 25,
		GCGraceHours:        72,
	}
	if issues := ValidateLifecyclePolicy(cfg); len(issues) != 0 {
		t.Fatalf("ValidateLifecyclePolicy returned issues: %v", issues)
	}
}

func TestValidateLifecyclePolicyRejectsUnsafePolicy(t *testing.T) {
	cfg := LifecyclePolicy{
		Name:                "unsafe",
		TransitionAfterDays: 30,
		ExpireAfterDays:     20,
		DeleteGraceDays:     0,
		LegalHoldAware:      false,
		DryRunRequired:      false,
		CompactionBudgetPct: 5,
		GCGraceHours:        0,
	}
	if issues := ValidateLifecyclePolicy(cfg); len(issues) < 5 {
		t.Fatalf("ValidateLifecyclePolicy returned too few issues: %v", issues)
	}
}
