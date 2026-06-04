package main

import "testing"

func TestScoreAnswerHealthy(t *testing.T) {
	card := AnswerScorecard{
		Clarifications:   3,
		HasSizing:        true,
		HasRequirements:  true,
		HasArchitecture:  true,
		DeepDives:        2,
		HasFailureModes:  true,
		HasObservability: true,
		HasTradeoffs:     true,
		HasRollout:       true,
	}
	if gaps := ScoreAnswer(card); len(gaps) != 0 {
		t.Fatalf("ScoreAnswer returned gaps: %v", gaps)
	}
}

func TestScoreAnswerFlagsMissingCorePieces(t *testing.T) {
	card := AnswerScorecard{}
	if gaps := ScoreAnswer(card); len(gaps) < 7 {
		t.Fatalf("ScoreAnswer returned too few gaps: %v", gaps)
	}
}
