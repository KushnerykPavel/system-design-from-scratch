package main

import "testing"

func TestValidateRedirectPlanHealthy(t *testing.T) {
	plan := RedirectPlan{
		Name:                  "healthy",
		CodeLength:            8,
		CacheTTLSeconds:       300,
		UsesEdgeCache:         true,
		HasIdempotencyKeys:    true,
		AsyncAnalytics:        true,
		SupportsCustomAlias:   true,
		HotKeyMitigation:      true,
		PolicyScanBeforeWrite: true,
	}
	if issues := ValidateRedirectPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateRedirectPlan returned issues: %v", issues)
	}
}

func TestValidateRedirectPlanRejectsWeakPlan(t *testing.T) {
	plan := RedirectPlan{
		Name:                "weak",
		CodeLength:          4,
		CacheTTLSeconds:     0,
		SupportsCustomAlias: true,
	}
	if issues := ValidateRedirectPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateRedirectPlan returned too few issues: %v", issues)
	}
}
