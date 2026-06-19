package main

import (
	"testing"
)

func TestStrongHire(t *testing.T) {
	card := StripeConnectScorecard{
		Idempotency:     95,
		LedgerDesign:    90,
		APIDesign:       90,
		FraudHooks:      90,
		FailureRecovery: 90,
		Observability:   85,
	}
	result := EvaluateMock(card)
	if result.HireSignal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if result.WeightedScore < 85 {
		t.Fatalf("expected weighted score >= 85, got %d", result.WeightedScore)
	}
	if len(result.WeakAreas) != 0 {
		t.Fatalf("expected no weak areas, got %v", result.WeakAreas)
	}
}

func TestHireSignal(t *testing.T) {
	card := StripeConnectScorecard{
		Idempotency:     80,
		LedgerDesign:    75,
		APIDesign:       70,
		FraudHooks:      65,
		FailureRecovery: 70,
		Observability:   60,
	}
	result := EvaluateMock(card)
	if result.HireSignal != "hire" {
		t.Fatalf("expected hire, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
}

func TestNoHire(t *testing.T) {
	card := StripeConnectScorecard{
		Idempotency:     20,
		LedgerDesign:    15,
		APIDesign:       30,
		FraudHooks:      10,
		FailureRecovery: 20,
		Observability:   10,
	}
	result := EvaluateMock(card)
	if result.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if result.WeightedScore >= 45 {
		t.Fatalf("expected score < 45, got %d", result.WeightedScore)
	}
}

func TestWeakAreaIdentification(t *testing.T) {
	card := StripeConnectScorecard{
		Idempotency:     90,
		LedgerDesign:    85,
		APIDesign:       70,
		FraudHooks:      40, // weak
		FailureRecovery: 75,
		Observability:   35, // weak
	}
	result := EvaluateMock(card)

	weakSet := make(map[string]bool)
	for _, area := range result.WeakAreas {
		weakSet[area] = true
	}

	if !weakSet["FraudHooks"] {
		t.Errorf("expected FraudHooks in weak areas, got %v", result.WeakAreas)
	}
	if !weakSet["Observability"] {
		t.Errorf("expected Observability in weak areas, got %v", result.WeakAreas)
	}
	if weakSet["Idempotency"] {
		t.Errorf("Idempotency should not be a weak area (score=90), got %v", result.WeakAreas)
	}
}

func TestWeightedScoringAccuracy(t *testing.T) {
	// All scores = 100 → weighted score must be 100.
	perfect := StripeConnectScorecard{
		Idempotency:     100,
		LedgerDesign:    100,
		APIDesign:       100,
		FraudHooks:      100,
		FailureRecovery: 100,
		Observability:   100,
	}
	result := EvaluateMock(perfect)
	if result.WeightedScore != 100 {
		t.Fatalf("perfect scorecard should yield 100, got %d", result.WeightedScore)
	}
}

func TestMixedSignal(t *testing.T) {
	card := StripeConnectScorecard{
		Idempotency:     60,
		LedgerDesign:    55,
		APIDesign:       50,
		FraudHooks:      50,
		FailureRecovery: 50,
		Observability:   45,
	}
	result := EvaluateMock(card)
	if result.HireSignal != "mixed" {
		t.Fatalf("expected mixed, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
}
