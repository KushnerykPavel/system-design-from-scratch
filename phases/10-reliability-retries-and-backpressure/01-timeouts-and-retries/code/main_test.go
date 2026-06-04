package main

import "testing"

func TestAssessPolicyFlagsRetryStormRisk(t *testing.T) {
	got := AssessPolicy(RetryProfile{
		Fanout:              5,
		TimeoutMS:           80,
		P99LatencyMS:        180,
		MaxAttempts:         4,
		HasJitter:           false,
		PropagatesDeadline:  false,
		OperationIdempotent: true,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessPolicyKeepsSafeProfileLowRisk(t *testing.T) {
	got := AssessPolicy(RetryProfile{
		Fanout:              1,
		TimeoutMS:           220,
		P99LatencyMS:        180,
		MaxAttempts:         2,
		HasJitter:           true,
		PropagatesDeadline:  true,
		OperationIdempotent: true,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}

func TestAssessPolicyProtectsNonIdempotentOperation(t *testing.T) {
	got := AssessPolicy(RetryProfile{
		Fanout:              1,
		TimeoutMS:           200,
		P99LatencyMS:        150,
		MaxAttempts:         2,
		HasJitter:           true,
		PropagatesDeadline:  true,
		OperationIdempotent: false,
	})

	if got.Risk == "low" {
		t.Fatalf("risk = %q, want medium or high", got.Risk)
	}
}
