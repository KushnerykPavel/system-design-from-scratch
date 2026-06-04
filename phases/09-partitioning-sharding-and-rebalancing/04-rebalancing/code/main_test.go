package main

import "testing"

func TestPlanMovesRespectsBatchLimits(t *testing.T) {
	moves := []Move{
		{ID: "a", SizeGB: 100, BandwidthGb: 2, Risk: 2},
		{ID: "b", SizeGB: 100, BandwidthGb: 2, Risk: 2},
		{ID: "c", SizeGB: 100, BandwidthGb: 3, Risk: 3},
	}

	plan := planMoves(moves, 2, 4, 4)
	if len(plan.Batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(plan.Batches))
	}
	if len(plan.Batches[0]) != 2 {
		t.Fatalf("expected first batch to contain 2 moves, got %d", len(plan.Batches[0]))
	}
}

func TestPlanMovesSplitsWhenRiskTooHigh(t *testing.T) {
	moves := []Move{
		{ID: "a", SizeGB: 100, BandwidthGb: 1, Risk: 3},
		{ID: "b", SizeGB: 100, BandwidthGb: 1, Risk: 3},
	}

	plan := planMoves(moves, 2, 4, 4)
	if len(plan.Batches) != 2 {
		t.Fatalf("expected risk limit to split batches, got %d", len(plan.Batches))
	}
}
