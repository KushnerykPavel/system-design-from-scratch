package main

import "testing"

func TestValidateAssumptionLogAcceptsImpactfulAssumptions(t *testing.T) {
	t.Parallel()

	log := AssumptionLog{
		Assumptions: []Assumption{
			{Statement: "single-region v1", Category: "geography", Impact: "simplifies write path", Reversible: false},
			{Statement: "bounded eventual consistency is acceptable", Category: "consistency", Impact: "allows async replication", Reversible: true},
		},
	}

	if issues := ValidateAssumptionLog(log); len(issues) != 0 {
		t.Fatalf("ValidateAssumptionLog() returned issues: %v", issues)
	}
}

func TestValidateAssumptionLogRejectsMissingImpact(t *testing.T) {
	t.Parallel()

	log := AssumptionLog{
		Assumptions: []Assumption{
			{Statement: "high scale", Category: "scale"},
		},
	}

	if issues := ValidateAssumptionLog(log); len(issues) == 0 {
		t.Fatal("ValidateAssumptionLog() returned no issues for missing impact")
	}
}

func TestCountExpensiveReversals(t *testing.T) {
	t.Parallel()

	log := AssumptionLog{
		Assumptions: []Assumption{
			{Statement: "single region", Category: "geography", Impact: "simple topology", Reversible: false},
			{Statement: "cheap storage class", Category: "scale", Impact: "lower cost", Reversible: true},
		},
	}

	if got, want := CountExpensiveReversals(log), 1; got != want {
		t.Fatalf("CountExpensiveReversals() = %d, want %d", got, want)
	}
}
