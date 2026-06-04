package main

import "testing"

func TestNextStepsOrdersMigrationGates(t *testing.T) {
	steps := nextSteps([]Cohort{
		{ID: "a", BackfillDone: false},
		{ID: "b", BackfillDone: true, DualWriteReady: false},
		{ID: "c", BackfillDone: true, DualWriteReady: true, ClientCompatible: false},
		{ID: "d", BackfillDone: true, DualWriteReady: true, ClientCompatible: true, ParityMismatch: 2},
		{ID: "e", BackfillDone: true, DualWriteReady: true, ClientCompatible: true, ParityMismatch: 0},
	})

	expected := []string{
		"start_backfill",
		"enable_dual_write",
		"hold_for_client_compat",
		"investigate_parity",
		"cutover",
	}

	for i, step := range steps {
		if step.Action != expected[i] {
			t.Fatalf("step %d: expected %s, got %s", i, expected[i], step.Action)
		}
	}
}
