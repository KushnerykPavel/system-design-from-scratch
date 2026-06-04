package main

import "testing"

func TestAssessBudgetHealthy(t *testing.T) {
	got := AssessBudget(BudgetInput{
		TotalEvents: 1000000,
		BadEvents:   200,
		TargetRatio: 0.999,
		WindowDays:  30,
		ElapsedDays: 10,
	})

	if got.Status != "healthy" {
		t.Fatalf("status = %q, want healthy", got.Status)
	}
}

func TestAssessBudgetWarningOnFastBurn(t *testing.T) {
	got := AssessBudget(BudgetInput{
		TotalEvents: 1000000,
		BadEvents:   500,
		TargetRatio: 0.999,
		WindowDays:  30,
		ElapsedDays: 5,
	})

	if got.Status != "warning" {
		t.Fatalf("status = %q, want warning", got.Status)
	}
}

func TestAssessBudgetCriticalWhenBudgetExhausted(t *testing.T) {
	got := AssessBudget(BudgetInput{
		TotalEvents: 1000000,
		BadEvents:   1500,
		TargetRatio: 0.999,
		WindowDays:  30,
		ElapsedDays: 20,
	})

	if got.Status != "critical" {
		t.Fatalf("status = %q, want critical", got.Status)
	}
}
