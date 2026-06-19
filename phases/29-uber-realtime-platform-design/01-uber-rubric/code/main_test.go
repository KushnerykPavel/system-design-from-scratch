package main

import "testing"

func TestEvaluateRubricAllObserved(t *testing.T) {
	signals := DefaultSignals()
	for i := range signals {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.ScorePercent != 100 {
		t.Fatalf("expected score 100, got %d", r.ScorePercent)
	}
	if r.HireSignal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %s", r.HireSignal)
	}
	if len(r.MissingAreas) != 0 {
		t.Fatalf("expected no missing areas, got %v", r.MissingAreas)
	}
}

func TestEvaluateRubricNoneObserved(t *testing.T) {
	signals := DefaultSignals()
	r := EvaluateRubric(signals)
	if r.ScorePercent != 0 {
		t.Fatalf("expected score 0, got %d", r.ScorePercent)
	}
	if r.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire, got %s", r.HireSignal)
	}
	if len(r.MissingAreas) != len(signals) {
		t.Fatalf("expected %d missing areas, got %d", len(signals), len(r.MissingAreas))
	}
}

func TestEvaluateRubricThresholdHire(t *testing.T) {
	// 6 of 8 signals = 75% → "hire"
	signals := DefaultSignals()
	for i := 0; i < 6; i++ {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.HireSignal != "hire" {
		t.Fatalf("expected hire at 75%%, got %s (score=%d)", r.HireSignal, r.ScorePercent)
	}
}

func TestEvaluateRubricThresholdMixed(t *testing.T) {
	// 3 of 8 signals = 37% → "mixed" (threshold is >=38, so 37% is no-hire)
	// Use 4 of 8 = 50% → "mixed"
	signals := DefaultSignals()
	for i := 0; i < 4; i++ {
		signals[i].Observed = true
	}
	r := EvaluateRubric(signals)
	if r.HireSignal != "mixed" {
		t.Fatalf("expected mixed at 50%%, got %s (score=%d)", r.HireSignal, r.ScorePercent)
	}
}

func TestEvaluateRubricMissingSurgePricing(t *testing.T) {
	// Main scenario: all observed except surge-pricing-design → 7/8 = 87% → "hire" (just below 88%)
	signals := DefaultSignals()
	for i := range signals {
		if signals[i].Category != "surge-pricing-design" {
			signals[i].Observed = true
		}
	}
	r := EvaluateRubric(signals)
	if r.ObservedSignals != 7 {
		t.Fatalf("expected 7 observed, got %d", r.ObservedSignals)
	}
	// 7/8 = 87%, which is below the 88% strong-hire threshold
	if r.HireSignal == "strong-hire" {
		t.Fatalf("expected hire (not strong-hire) when surge pricing missing, got %s", r.HireSignal)
	}
	found := false
	for _, m := range r.MissingAreas {
		if len(m) > 0 && m[:len("surge-pricing-design")] == "surge-pricing-design" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected surge-pricing-design in missing areas, got %v", r.MissingAreas)
	}
}
