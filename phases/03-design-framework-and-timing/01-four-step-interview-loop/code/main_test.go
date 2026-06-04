package main

import "testing"

func TestDefaultPlanFitsInterviewWindow(t *testing.T) {
	plan := DefaultPlan("Design a rate limiter")
	if got := TotalMinutes(plan); got != 45 {
		t.Fatalf("TotalMinutes() = %d, want 45", got)
	}
	if issues := ValidatePlan(plan); len(issues) != 0 {
		t.Fatalf("ValidatePlan() returned issues: %v", issues)
	}
}

func TestValidatePlanCatchesOverBudget(t *testing.T) {
	plan := Plan{
		Prompt: "Design chat",
		Stages: []Stage{
			{Name: "clarify", Minutes: 10},
			{Name: "size", Minutes: 10},
			{Name: "high_level_design", Minutes: 10},
			{Name: "deep_dive", Minutes: 10},
			{Name: "wrap_up", Minutes: 10},
		},
	}
	if issues := ValidatePlan(plan); len(issues) == 0 {
		t.Fatal("ValidatePlan() returned no issues for an over-budget plan")
	}
}
