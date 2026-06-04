package main

import "testing"

func TestValidateGatewayConfigAcceptsHealthyConfig(t *testing.T) {
	cfg := GatewayConfig{
		PublicHost:            true,
		TLSMode:               "terminate",
		RetryBudgetPercent:    10,
		OriginShieldRegion:    "us-east",
		PolicyVersion:         "v3",
		AuthDependencyMode:    "hybrid",
		ObservabilityComplete: true,
	}
	if issues := ValidateGatewayConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateGatewayConfig() returned issues: %v", issues)
	}
}

func TestValidateGatewayConfigRejectsMissingGuardrails(t *testing.T) {
	cfg := GatewayConfig{
		PublicHost:            true,
		TLSMode:               "",
		RetryBudgetPercent:    0,
		OriginShieldRegion:    "",
		PolicyVersion:         "",
		AuthDependencyMode:    "none",
		ObservabilityComplete: false,
	}
	if issues := ValidateGatewayConfig(cfg); len(issues) < 4 {
		t.Fatalf("ValidateGatewayConfig() returned too few issues: %v", issues)
	}
}
