package main

import "testing"

func TestEvaluateRubricAllObserved(t *testing.T) {
	signals := DefaultSignals()
	for i := range signals {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.ScorePercent != 100 {
		t.Fatalf("expected 100, got %d", r.ScorePercent)
	}
	if r.HireSignal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %s", r.HireSignal)
	}
}

func TestEvaluateRubricNoneObserved(t *testing.T) {
	signals := DefaultSignals()
	r := EvaluateRubric(signals)
	if r.ScorePercent != 0 {
		t.Fatalf("expected 0, got %d", r.ScorePercent)
	}
	if r.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire, got %s", r.HireSignal)
	}
}

func TestEvaluateRubricMixed(t *testing.T) {
	signals := DefaultSignals()
	// 4 out of 8 = 50%, which falls in the mixed band (>= 38 and < 63).
	for i := 0; i < 4; i++ {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.HireSignal != "mixed" {
		t.Fatalf("expected mixed, got %s (score=%d)", r.HireSignal, r.ScorePercent)
	}
}
