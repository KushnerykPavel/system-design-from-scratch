package main

import "testing"

func TestValidateSessionPlan(t *testing.T) {
	t.Parallel()

	plan := SessionPlan{
		DurationMinutes: 45,
		Stages: []Stage{
			{Name: "pre_brief", Minutes: 3},
			{Name: "live", Minutes: 28},
			{Name: "feedback", Minutes: 9},
			{Name: "debrief", Minutes: 5},
		},
	}

	if issues := ValidateSessionPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateSessionPlan() returned issues: %v", issues)
	}
}

func TestValidateSessionPlanRejectsMissingStages(t *testing.T) {
	t.Parallel()

	plan := SessionPlan{
		DurationMinutes: 20,
		Stages: []Stage{
			{Name: "live", Minutes: 18},
		},
	}

	if issues := ValidateSessionPlan(plan); len(issues) == 0 {
		t.Fatal("ValidateSessionPlan() returned no issues for a broken plan")
	}
}

func TestTotalMinutes(t *testing.T) {
	t.Parallel()

	plan := SessionPlan{
		Stages: []Stage{
			{Name: "pre_brief", Minutes: 2},
			{Name: "live", Minutes: 10},
			{Name: "feedback", Minutes: 4},
		},
	}

	if got := TotalMinutes(plan); got != 16 {
		t.Fatalf("TotalMinutes() = %d, want 16", got)
	}
}
