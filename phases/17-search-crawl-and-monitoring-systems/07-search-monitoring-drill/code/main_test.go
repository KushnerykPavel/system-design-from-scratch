package main

import "testing"

func TestValidateDrillScorecardHealthy(t *testing.T) {
	card := DrillScorecard{
		Name:            "healthy",
		ClarifiedScope:  true,
		SizedWorkload:   true,
		ChosenDeepDive:  true,
		FailureModes:    true,
		Observability:   true,
		Tradeoffs:       true,
		RedesignHandled: true,
	}
	if issues := ValidateDrillScorecard(card); len(issues) != 0 {
		t.Fatalf("ValidateDrillScorecard returned issues: %v", issues)
	}
}

func TestValidateDrillScorecardWeak(t *testing.T) {
	card := DrillScorecard{Name: "weak"}
	if issues := ValidateDrillScorecard(card); len(issues) < 7 {
		t.Fatalf("ValidateDrillScorecard returned too few issues: %v", issues)
	}
}
