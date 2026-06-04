package main

import "testing"

func TestValidatePolicyAcceptsHealthyPolicy(t *testing.T) {
	policy := SteeringPolicy{
		Name:                   "good",
		IngressMode:            "anycast",
		PropagationSeconds:     30,
		HealthSignalWindowSecs: 20,
		HasCapacitySignals:     true,
		SupportsAffinity:       true,
		HasEmergencyOverride:   true,
		HasRouteExplanation:    true,
	}
	if issues := ValidatePolicy(policy); len(issues) != 0 {
		t.Fatalf("ValidatePolicy() returned issues: %v", issues)
	}
}

func TestValidatePolicyRejectsWeakPolicy(t *testing.T) {
	policy := SteeringPolicy{
		Name:                   "bad",
		IngressMode:            "magic",
		PropagationSeconds:     600,
		HealthSignalWindowSecs: 5,
		HasCapacitySignals:     false,
		SupportsAffinity:       false,
		HasEmergencyOverride:   false,
		HasRouteExplanation:    false,
	}
	if issues := ValidatePolicy(policy); len(issues) < 4 {
		t.Fatalf("ValidatePolicy() returned too few issues: %v", issues)
	}
}
