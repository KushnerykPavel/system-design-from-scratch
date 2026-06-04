package main

import "testing"

func TestMissingIngredientsFlagsIncompleteReframing(t *testing.T) {
	t.Parallel()

	prompt := ReframedPrompt{
		System:       "url shortener",
		CoreWorkflow: "create and resolve short links",
	}

	if got := len(MissingIngredients(prompt)); got != 4 {
		t.Fatalf("len(MissingIngredients()) = %d, want 4", got)
	}
}

func TestScoreReframingReturnsFullScore(t *testing.T) {
	t.Parallel()

	prompt := ReframedPrompt{
		System:        "photo sharing",
		CoreWorkflow:  "upload and read photos",
		ScopeCut:      "no collaboration",
		Priority:      "read latency",
		Assumption:    "single-region v1",
		WorkloadShape: "read-heavy with hotspot bursts",
	}

	if got, want := ScoreReframing(prompt), 6; got != want {
		t.Fatalf("ScoreReframing() = %d, want %d", got, want)
	}
}
