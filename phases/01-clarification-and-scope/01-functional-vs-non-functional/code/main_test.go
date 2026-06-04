package main

import "testing"

func TestValidateRequirementSetAcceptsMixedRequirements(t *testing.T) {
	t.Parallel()

	set := RequirementSet{
		Requirements: []Requirement{
			{Text: "users can upload files", Kind: "functional", Priority: 1},
			{Text: "p99 under 250 ms", Kind: "non_functional", Priority: 1, Driver: true},
			{Text: "survive zone loss", Kind: "non_functional", Priority: 2},
		},
	}

	if issues := ValidateRequirementSet(set); len(issues) != 0 {
		t.Fatalf("ValidateRequirementSet() returned issues: %v", issues)
	}
}

func TestValidateRequirementSetRejectsMissingDriver(t *testing.T) {
	t.Parallel()

	set := RequirementSet{
		Requirements: []Requirement{
			{Text: "users can search", Kind: "functional", Priority: 1},
			{Text: "cheap to operate", Kind: "non_functional", Priority: 1},
		},
	}

	if issues := ValidateRequirementSet(set); len(issues) == 0 {
		t.Fatal("ValidateRequirementSet() returned no issues for missing dominant driver")
	}
}

func TestRankedNonFunctionalSortsByPriority(t *testing.T) {
	t.Parallel()

	set := RequirementSet{
		Requirements: []Requirement{
			{Text: "availability", Kind: "non_functional", Priority: 2},
			{Text: "latency", Kind: "non_functional", Priority: 1},
			{Text: "upload file", Kind: "functional", Priority: 1},
		},
	}

	ranked := RankedNonFunctional(set)
	if got, want := ranked[0].Text, "latency"; got != want {
		t.Fatalf("RankedNonFunctional()[0] = %q, want %q", got, want)
	}
}
