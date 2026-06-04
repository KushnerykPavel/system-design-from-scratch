package main

import "testing"

func TestScoreAnswerCountsCoverage(t *testing.T) {
	got := ScoreAnswer(DrillAnswer{
		HasTimeoutOwnership:  true,
		HasIdempotentIngest:  true,
		HasRetryBudget:       false,
		HasLoadShedding:      true,
		HasAsyncBackpressure: false,
		HasIsolation:         true,
		HasObservability:     true,
	})

	if got.Score != 5 {
		t.Fatalf("score = %d, want 5", got.Score)
	}
	if len(got.Missing) != 2 {
		t.Fatalf("missing = %d items, want 2", len(got.Missing))
	}
}

func TestScoreAnswerFullCoverage(t *testing.T) {
	got := ScoreAnswer(DrillAnswer{
		HasTimeoutOwnership:  true,
		HasIdempotentIngest:  true,
		HasRetryBudget:       true,
		HasLoadShedding:      true,
		HasAsyncBackpressure: true,
		HasIsolation:         true,
		HasObservability:     true,
	})

	if got.Score != 7 {
		t.Fatalf("score = %d, want 7", got.Score)
	}
	if len(got.Missing) != 0 {
		t.Fatalf("missing = %v, want none", got.Missing)
	}
}
