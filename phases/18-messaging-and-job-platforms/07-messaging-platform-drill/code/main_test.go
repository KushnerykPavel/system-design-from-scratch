package main

import "testing"

func TestValidateDrillScorecard(t *testing.T) {
	full := DrillScorecard{
		ClarifiedScope:       true,
		SizedWorkload:        true,
		ExplainedGuarantees:  true,
		ChosenDeepDive:       true,
		CoveredFailures:      true,
		CoveredObservability: true,
		HandledRedesign:      true,
	}
	if issues := ValidateDrillScorecard(full); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	full.HandledRedesign = false
	full.ExplainedGuarantees = false
	if issues := ValidateDrillScorecard(full); len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
}
