package main

import "testing"

func TestTotalAndOverBudget(t *testing.T) {
	plan := Default45MinutePlan()
	if got := Total(plan); got != 45 {
		t.Fatalf("expected 45 minutes, got %d", got)
	}
	if got := OverBudget(Total(plan), 45); got != 0 {
		t.Fatalf("expected no overage, got %d", got)
	}

	over := []Segment{
		{Name: "clarify", Minutes: 10},
		{Name: "size", Minutes: 8},
		{Name: "design", Minutes: 15},
		{Name: "deep_dive", Minutes: 15},
	}
	if got := OverBudget(Total(over), 45); got != 3 {
		t.Fatalf("expected 3 minutes over budget, got %d", got)
	}
}
