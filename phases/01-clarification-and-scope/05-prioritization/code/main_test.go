package main

import "testing"

func TestValidatePrioritySetAcceptsRankedRequirements(t *testing.T) {
	t.Parallel()

	set := PrioritySet{
		Priorities: []Priority{
			{Requirement: "latency", Rank: 1, Rationale: "user-facing reads dominate"},
			{Requirement: "availability", Rank: 2, Rationale: "global access matters"},
			{Requirement: "cost", Rank: 3, Rationale: "optimize after correctness"},
		},
	}

	if issues := ValidatePrioritySet(set); len(issues) != 0 {
		t.Fatalf("ValidatePrioritySet() returned issues: %v", issues)
	}
}

func TestValidatePrioritySetRejectsDuplicateRanks(t *testing.T) {
	t.Parallel()

	set := PrioritySet{
		Priorities: []Priority{
			{Requirement: "latency", Rank: 1, Rationale: "fast"},
			{Requirement: "cost", Rank: 1, Rationale: "cheap"},
		},
	}

	if issues := ValidatePrioritySet(set); len(issues) == 0 {
		t.Fatal("ValidatePrioritySet() returned no issues for duplicate ranks")
	}
}

func TestSortedPriorities(t *testing.T) {
	t.Parallel()

	set := PrioritySet{
		Priorities: []Priority{
			{Requirement: "availability", Rank: 2},
			{Requirement: "latency", Rank: 1},
		},
	}

	sorted := SortedPriorities(set)
	if got, want := sorted[0].Requirement, "latency"; got != want {
		t.Fatalf("SortedPriorities()[0] = %q, want %q", got, want)
	}
}
