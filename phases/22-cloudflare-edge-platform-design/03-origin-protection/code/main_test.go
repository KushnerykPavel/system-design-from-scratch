package main

import "testing"

func TestValidateFailoverConfigAcceptsReasonableDefaults(t *testing.T) {
	cfg := FailoverConfig{
		PoolName:             "pool-a",
		HealthMode:           "hybrid",
		RetryBudgetPercent:   10,
		CooldownSeconds:      20,
		ShieldEnabled:        true,
		CrossRegionFailover:  true,
		Explainability:       true,
		ReserveCapacityRatio: 150,
	}
	if issues := ValidateFailoverConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateFailoverConfig() returned issues: %v", issues)
	}
}

func TestValidateFailoverConfigRejectsWeakGuardrails(t *testing.T) {
	cfg := FailoverConfig{
		PoolName:             "",
		HealthMode:           "none",
		RetryBudgetPercent:   0,
		CooldownSeconds:      0,
		ShieldEnabled:        false,
		CrossRegionFailover:  true,
		Explainability:       false,
		ReserveCapacityRatio: 100,
	}
	if issues := ValidateFailoverConfig(cfg); len(issues) < 5 {
		t.Fatalf("ValidateFailoverConfig() returned too few issues: %v", issues)
	}
}
