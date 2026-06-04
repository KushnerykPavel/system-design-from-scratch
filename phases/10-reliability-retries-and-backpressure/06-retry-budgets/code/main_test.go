package main

import "testing"

func TestAssessBudgetPolicyFlagsUnsafeHedging(t *testing.T) {
	got := AssessBudgetPolicy(BudgetPolicy{
		BaseQPS:              180000,
		RetryBudgetRatio:     0.08,
		HedgeBudgetRatio:     0.08,
		HedgeAfterMS:         0,
		SupportsCancellation: false,
		SafeToHedge:          false,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessBudgetPolicyApprovesBoundedSpeculation(t *testing.T) {
	got := AssessBudgetPolicy(BudgetPolicy{
		BaseQPS:              180000,
		RetryBudgetRatio:     0.05,
		HedgeBudgetRatio:     0.02,
		HedgeAfterMS:         120,
		SupportsCancellation: true,
		SafeToHedge:          true,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
