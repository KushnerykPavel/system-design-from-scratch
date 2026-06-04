package main

import "testing"

func TestValidateIndexUpdatePlanHealthy(t *testing.T) {
	plan := IndexUpdatePlan{
		Name:                    "healthy",
		DocumentLagSeconds:      300,
		RankingLagSeconds:       600,
		BlueGreenEnabled:        true,
		DualReadValidation:      true,
		DeleteTombstonesEnabled: true,
		BackfillWorkerPools:     1,
		SchemaCompatChecks:      true,
	}
	if issues := ValidateIndexUpdatePlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateIndexUpdatePlan returned issues: %v", issues)
	}
}

func TestValidateIndexUpdatePlanWeak(t *testing.T) {
	plan := IndexUpdatePlan{Name: "weak"}
	if issues := ValidateIndexUpdatePlan(plan); len(issues) < 6 {
		t.Fatalf("ValidateIndexUpdatePlan returned too few issues: %v", issues)
	}
}
