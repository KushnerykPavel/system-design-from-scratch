package main

import "testing"

func TestCoverageScoreAndStrength(t *testing.T) {
	summary := WrapUp{
		Risks:          2,
		TradeOffs:      2,
		Observability:  true,
		RolloutMention: true,
	}
	if got := CoverageScore(summary); got != 6 {
		t.Fatalf("expected score 6, got %d", got)
	}
	if !StrongWrapUp(summary) {
		t.Fatal("expected strong wrap-up")
	}
}
