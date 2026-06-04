package main

import "testing"

func TestTotalCostAndMostExpensiveItem(t *testing.T) {
	items := []CostItem{
		{Name: "compute", UnitCost: 200, Units: 180},
		{Name: "storage", UnitCost: 25, Units: 120},
		{Name: "egress", UnitCost: 0.02, Units: 8000000},
	}

	if got := TotalCost(items); got != 199000 {
		t.Fatalf("unexpected total cost: %.2f", got)
	}
	if got := MostExpensiveItem(items); got != "egress" {
		t.Fatalf("unexpected most expensive item: %s", got)
	}
}
