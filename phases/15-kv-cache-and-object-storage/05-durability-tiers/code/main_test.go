package main

import "testing"

func TestValidateDurabilityTierAcceptsReplicatedTier(t *testing.T) {
	cfg := DurabilityTier{
		Name:           "standard",
		ReplicaCopies:  3,
		MaxRepairHours: 6,
		RestoreTested:  true,
	}
	if issues := ValidateDurabilityTier(cfg); len(issues) != 0 {
		t.Fatalf("ValidateDurabilityTier returned issues: %v", issues)
	}
}

func TestValidateDurabilityTierRejectsWeakTier(t *testing.T) {
	cfg := DurabilityTier{
		Name:           "weak",
		ReplicaCopies:  2,
		GeoRedundant:   false,
		MaxRepairHours: 0,
		RestoreTested:  false,
	}
	if issues := ValidateDurabilityTier(cfg); len(issues) < 3 {
		t.Fatalf("ValidateDurabilityTier returned too few issues: %v", issues)
	}
}
