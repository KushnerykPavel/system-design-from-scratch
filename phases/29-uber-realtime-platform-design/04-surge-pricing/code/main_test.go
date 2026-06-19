package main

import "testing"

// TestNoSurge verifies that ratio <= 1.0 produces a 1.0x multiplier.
func TestNoSurge(t *testing.T) {
	zone := Zone{ID: "quiet", AvailableDrivers: 20, OpenRequests: 10}
	got := SurgeMultiplier(zone)
	if got != 1.0 {
		t.Errorf("expected 1.0x for ratio 0.5, got %.1f", got)
	}
}

// TestEqualSupplyDemand verifies that ratio exactly 1.0 produces 1.2x (first surge tier).
func TestEqualSupplyDemand(t *testing.T) {
	zone := Zone{ID: "balanced", AvailableDrivers: 5, OpenRequests: 5}
	got := SurgeMultiplier(zone)
	if got != 1.2 {
		t.Errorf("expected 1.2x for ratio 1.0, got %.1f", got)
	}
}

// TestModerateSurge verifies that ratio 1.5 maps to 2.0x multiplier.
func TestModerateSurge(t *testing.T) {
	zone := Zone{ID: "moderate", AvailableDrivers: 2, OpenRequests: 3}
	// ratio = 3/2 = 1.5 → should be 2.0x
	got := SurgeMultiplier(zone)
	if got != 2.0 {
		t.Errorf("expected 2.0x for ratio 1.5, got %.1f", got)
	}
}

// TestHighSurge verifies that ratio in 2.0–4.0 range maps to 4.0x.
func TestHighSurge(t *testing.T) {
	zone := Zone{ID: "high", AvailableDrivers: 1, OpenRequests: 3}
	// ratio = 3.0 → 4.0x
	got := SurgeMultiplier(zone)
	if got != 4.0 {
		t.Errorf("expected 4.0x for ratio 3.0, got %.1f", got)
	}
}

// TestCapAt8x verifies that very high demand ratios are capped at 8.0x.
func TestCapAt8x(t *testing.T) {
	zone := Zone{ID: "disaster", AvailableDrivers: 1, OpenRequests: 100}
	got := SurgeMultiplier(zone)
	if got != 8.0 {
		t.Errorf("expected 8.0x cap for ratio 100, got %.1f", got)
	}
}

// TestZeroDriversCappedAt8x verifies that zero supply is treated as supply=1.
func TestZeroDriversCappedAt8x(t *testing.T) {
	zone := Zone{ID: "no_drivers", AvailableDrivers: 0, OpenRequests: 10}
	got := SurgeMultiplier(zone)
	// ratio = 10/1 = 10 → capped at 8x
	if got != 8.0 {
		t.Errorf("expected 8.0x when AvailableDrivers=0, got %.1f", got)
	}
}

// TestConfirmationRequiredAt2x verifies that confirmation is required at exactly 2.0x.
func TestConfirmationRequiredAt2x(t *testing.T) {
	if !RequiresConfirmation(2.0) {
		t.Error("expected confirmation required at 2.0x")
	}
}

// TestNoConfirmationBelow2x verifies that confirmation is not required below 2.0x.
func TestNoConfirmationBelow2x(t *testing.T) {
	multipliers := []float64{1.0, 1.2, 1.5}
	for _, m := range multipliers {
		if RequiresConfirmation(m) {
			t.Errorf("expected no confirmation required at %.1fx", m)
		}
	}
}

// TestConfirmationRequiredAbove2x verifies high multipliers also require confirmation.
func TestConfirmationRequiredAbove2x(t *testing.T) {
	multipliers := []float64{2.0, 4.0, 8.0}
	for _, m := range multipliers {
		if !RequiresConfirmation(m) {
			t.Errorf("expected confirmation required at %.1fx", m)
		}
	}
}
