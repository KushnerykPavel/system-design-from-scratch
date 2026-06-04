package main

import "testing"

func TestDetectAntiPatternsFindsMissingBehaviors(t *testing.T) {
	t.Parallel()

	summary := PracticeSummary{}
	issues := DetectAntiPatterns(summary)
	if len(issues) != 5 {
		t.Fatalf("DetectAntiPatterns() returned %d issues, want 5", len(issues))
	}
}

func TestDetectAntiPatternsAcceptsHealthyOpening(t *testing.T) {
	t.Parallel()

	summary := PracticeSummary{
		ClarifiedScope:     true,
		RankedRequirements: true,
		LoggedAssumptions:  true,
		NamedWorkloadShape: true,
		StatedScopeCut:     true,
	}

	if issues := DetectAntiPatterns(summary); len(issues) != 0 {
		t.Fatalf("DetectAntiPatterns() returned issues: %v", issues)
	}
}
