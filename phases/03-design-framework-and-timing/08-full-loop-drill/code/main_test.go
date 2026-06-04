package main

import "testing"

func TestCoverage(t *testing.T) {
	answer := AnswerCoverage{
		Clarify:       true,
		Sizing:        true,
		HighLevel:     true,
		DeepDive:      true,
		FailureModes:  true,
		Observability: true,
		TradeOffs:     true,
		Redesign:      true,
	}
	if got := CoveredCount(answer); got != 8 {
		t.Fatalf("expected 8 covered areas, got %d", got)
	}
	if !IsFullLoop(answer) {
		t.Fatal("expected answer to cover the full loop")
	}
}
