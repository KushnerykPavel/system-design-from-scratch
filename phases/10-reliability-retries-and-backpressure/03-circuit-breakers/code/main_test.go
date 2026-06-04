package main

import "testing"

func TestAssessBreakerFlagsSharedFallback(t *testing.T) {
	got := AssessBreaker(BreakerConfig{
		TripErrorRate:        0.9,
		TripTimeoutRate:      0.9,
		UsesSaturationSignal: false,
		FallbackIndependent:  false,
		HalfOpenMaxProbes:    50,
		ScopedPerDependency:  false,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessBreakerApprovesScopedIndependentFallback(t *testing.T) {
	got := AssessBreaker(BreakerConfig{
		TripErrorRate:        0.35,
		TripTimeoutRate:      0.25,
		UsesSaturationSignal: true,
		FallbackIndependent:  true,
		HalfOpenMaxProbes:    3,
		ScopedPerDependency:  true,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
