package main

import "testing"

func TestValidateProximityConfigHealthy(t *testing.T) {
	cfg := ProximityConfig{
		Name:                 "healthy",
		DefaultRadiusM:       1000,
		MaxRadiusM:           10000,
		MaxCandidates:        400,
		HotCellReplicas:      3,
		LocationTTLSeconds:   15,
		AvailabilityCacheTTL: 5,
		ExactRerankEnabled:   true,
	}
	if issues := ValidateProximityConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateProximityConfig returned issues: %v", issues)
	}
}

func TestValidateProximityConfigWeak(t *testing.T) {
	cfg := ProximityConfig{
		Name:                 "weak",
		DefaultRadiusM:       0,
		MaxRadiusM:           500,
		MaxCandidates:        5,
		HotCellReplicas:      1,
		LocationTTLSeconds:   1000,
		AvailabilityCacheTTL: 120,
	}
	if issues := ValidateProximityConfig(cfg); len(issues) < 5 {
		t.Fatalf("ValidateProximityConfig returned too few issues: %v", issues)
	}
}
