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

func TestEvaluateRubricHireThreshold(t *testing.T) {
	signals := DefaultSignals()
	for i := 0; i < 6; i++ {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.HireSignal != "hire" {
		t.Fatalf("expected hire, got %s (score=%d)", r.HireSignal, r.ScorePercent)
	}
}
