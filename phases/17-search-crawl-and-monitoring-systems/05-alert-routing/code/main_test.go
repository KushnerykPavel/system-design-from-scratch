package main

import "testing"

func TestValidateAlertPolicyHealthy(t *testing.T) {
	policy := AlertPolicy{
		Name:                 "healthy",
		EscalationLevels:     1,
		DedupWindowMinutes:   5,
		AckTimeoutMinutes:    10,
		OwnershipRequired:    true,
		RunbookRequired:      true,
		ProviderFallback:     true,
		MaintenanceSupported: true,
	}
	if issues := ValidateAlertPolicy(policy); len(issues) != 0 {
		t.Fatalf("ValidateAlertPolicy returned issues: %v", issues)
	}
}

func TestValidateAlertPolicyWeak(t *testing.T) {
	policy := AlertPolicy{Name: "weak"}
	if issues := ValidateAlertPolicy(policy); len(issues) < 6 {
		t.Fatalf("ValidateAlertPolicy returned too few issues: %v", issues)
	}
}
