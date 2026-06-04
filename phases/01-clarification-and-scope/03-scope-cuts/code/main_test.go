package main

import "testing"

func TestValidateScopePlanAcceptsCredibleCut(t *testing.T) {
	t.Parallel()

	plan := ScopePlan{
		CoreWorkflows:         []string{"upload", "download"},
		DeferredFeatures:      []string{"collaborative editing"},
		Reason:                "preserve the storage and sync core while reducing breadth",
		PreservesPromptIntent: true,
	}

	if issues := ValidateScopePlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateScopePlan() returned issues: %v", issues)
	}
}

func TestValidateScopePlanRejectsEvasiveCut(t *testing.T) {
	t.Parallel()

	plan := ScopePlan{
		DeferredFeatures:      []string{"storage"},
		Reason:                "too hard otherwise",
		PreservesPromptIntent: false,
	}

	if issues := ValidateScopePlan(plan); len(issues) == 0 {
		t.Fatal("ValidateScopePlan() returned no issues for an evasive plan")
	}
}

func TestComplexityReduction(t *testing.T) {
	t.Parallel()

	plan := ScopePlan{
		CoreWorkflows:    []string{"write"},
		DeferredFeatures: []string{"analytics", "admin", "multi-region"},
	}

	if got, want := ComplexityReduction(plan), 2; got != want {
		t.Fatalf("ComplexityReduction() = %d, want %d", got, want)
	}
}
