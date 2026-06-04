package main

import "testing"

func TestAssessIdempotencyPolicyRejectsUnsafeOrdering(t *testing.T) {
	got := AssessIdempotencyPolicy(IdempotencyPolicy{
		DurableStore:            true,
		ReserveBeforeSideEffect: false,
		StoresRequestHash:       true,
		ReturnsStoredResponse:   true,
		TTLHours:                24,
	}, 12)

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessIdempotencyPolicyFlagsShortTTL(t *testing.T) {
	got := AssessIdempotencyPolicy(IdempotencyPolicy{
		DurableStore:            true,
		ReserveBeforeSideEffect: true,
		StoresRequestHash:       true,
		ReturnsStoredResponse:   false,
		TTLHours:                2,
	}, 12)

	if got.Risk == "low" {
		t.Fatalf("risk = %q, want medium or high", got.Risk)
	}
}

func TestAssessIdempotencyPolicyApprovesStrongPolicy(t *testing.T) {
	got := AssessIdempotencyPolicy(IdempotencyPolicy{
		DurableStore:            true,
		ReserveBeforeSideEffect: true,
		StoresRequestHash:       true,
		ReturnsStoredResponse:   true,
		TTLHours:                24,
	}, 12)

	if !got.Safe {
		t.Fatalf("safe = false, want true")
	}
}
