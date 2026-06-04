package main

import "testing"

func TestValidateConfigAcceptsReasonableDefaults(t *testing.T) {
	cfg := LimiterConfig{
		Name:               "default",
		Mode:               "token_bucket",
		Store:              "redis",
		LocalTTLMillis:     100,
		Burst:              20,
		RetryBudgetPercent: 5,
		FailOpen:           true,
		SharedEnforcement:  true,
		Explainability:     true,
	}
	if issues := ValidateConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateConfig() returned issues: %v", issues)
	}
}

func TestValidateConfigRejectsWeakSettings(t *testing.T) {
	cfg := LimiterConfig{
		Name:              "broken",
		Mode:              "counter",
		Store:             "",
		LocalTTLMillis:    0,
		Burst:             0,
		FailOpen:          true,
		SharedEnforcement: true,
		Explainability:    false,
	}
	if issues := ValidateConfig(cfg); len(issues) < 4 {
		t.Fatalf("ValidateConfig() returned too few issues: %v", issues)
	}
}
