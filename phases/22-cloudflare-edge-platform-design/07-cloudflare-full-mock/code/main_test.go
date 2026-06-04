package main

import "testing"

func TestValidateMockScorecardAcceptsCompleteAnswer(t *testing.T) {
	card := MockScorecard{
		ScopeClarity:            true,
		SizingIncluded:          true,
		ControlDataPlaneSplit:   true,
		DeepDiveIncluded:        true,
		FailureCoverage:         true,
		ObservabilityCoverage:   true,
		RedesignIncluded:        true,
		AmplificationRiskNamed:  true,
		ExplainabilityAddressed: true,
	}
	if issues := ValidateMockScorecard(card); len(issues) != 0 {
		t.Fatalf("ValidateMockScorecard() returned issues: %v", issues)
	}
}

func TestValidateMockScorecardRejectsMissingSignals(t *testing.T) {
	card := MockScorecard{}
	if issues := ValidateMockScorecard(card); len(issues) < 6 {
		t.Fatalf("ValidateMockScorecard() returned too few issues: %v", issues)
	}
}
