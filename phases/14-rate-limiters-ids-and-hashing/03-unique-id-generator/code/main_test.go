package main

import "testing"

func TestValidatePlanAcceptsReasonableSnowflakePlan(t *testing.T) {
	plan := IDPlan{
		Name:                   "good",
		Strategy:               "snowflake",
		PeakWritesPerSecond:    200000,
		RequiresSortability:    true,
		RequiresGuessability:   false,
		Regions:                3,
		ClientsGenerateIDs:     false,
		ClockDiscipline:        true,
		StorageLocalityMatters: true,
	}
	if issues := ValidatePlan(plan); len(issues) != 0 {
		t.Fatalf("ValidatePlan() returned issues: %v", issues)
	}
}

func TestValidatePlanRejectsWeakSettings(t *testing.T) {
	plan := IDPlan{
		Name:                   "bad",
		Strategy:               "random",
		PeakWritesPerSecond:    0,
		RequiresSortability:    true,
		RequiresGuessability:   true,
		Regions:                2,
		ClientsGenerateIDs:     false,
		ClockDiscipline:        false,
		StorageLocalityMatters: true,
	}
	if issues := ValidatePlan(plan); len(issues) < 3 {
		t.Fatalf("ValidatePlan() returned too few issues: %v", issues)
	}
}
