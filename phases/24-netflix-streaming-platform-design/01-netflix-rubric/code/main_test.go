package main

import "testing"

func TestEvaluateRubricAllObserved(t *testing.T) {
	signals := DefaultSignals()
	for i := range signals {
		signals[i].Observed = true
	}
	result := EvaluateRubric(signals)
	if result.ScorePercent != 100 {
		t.Fatalf("expected 100%% score, got %d", result.ScorePercent)
	}
	if result.HireSignal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %s", result.HireSignal)
	}
	if len(result.MissingAreas) != 0 {
		t.Fatalf("expected no missing areas, got %v", result.MissingAreas)
	}
}

func TestEvaluateRubricNoneObserved(t *testing.T) {
	signals := DefaultSignals()
	result := EvaluateRubric(signals)
	if result.ScorePercent != 0 {
		t.Fatalf("expected 0%% score, got %d", result.ScorePercent)
	}
	if result.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire, got %s", result.HireSignal)
	}
	if len(result.MissingAreas) != len(signals) {
		t.Fatalf("expected %d missing areas, got %d", len(signals), len(result.MissingAreas))
	}
}

func TestEvaluateRubricHireThreshold(t *testing.T) {
	signals := DefaultSignals()
	// Observe exactly 5 of 8 = 62% -> mixed, 6 of 8 = 75% -> hire
	for i := 0; i < 6; i++ {
		signals[i].Observed = true
	}
	result := EvaluateRubric(signals)
	if result.HireSignal != "hire" {
		t.Fatalf("expected hire at 75%%, got %s (score=%d)", result.HireSignal, result.ScorePercent)
	}
}
