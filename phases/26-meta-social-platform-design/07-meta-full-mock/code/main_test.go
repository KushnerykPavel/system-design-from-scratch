package main

import (
	"testing"
)

func TestScoreMockStrongHire(t *testing.T) {
	card := MockScorecard{
		Clarification: 100,
		Capacity:      100,
		Architecture:  100,
		DeepDive:      100,
		FailureModes:  100,
		Tradeoffs:     100,
		Observability: 100,
	}
	result := ScoreMock(card)

	if result.TotalScore != 100.0 {
		t.Errorf("expected total score 100.0, got %.2f", result.TotalScore)
	}
	if result.HireSignal != "strong-hire" {
		t.Errorf("expected strong-hire, got %s", result.HireSignal)
	}
	if len(result.WeakAreas) != 0 {
		t.Errorf("expected no weak areas for perfect score, got %v", result.WeakAreas)
	}
}

func TestScoreMockNoHire(t *testing.T) {
	card := MockScorecard{
		Clarification: 0,
		Capacity:      0,
		Architecture:  0,
		DeepDive:      0,
		FailureModes:  0,
		Tradeoffs:     0,
		Observability: 0,
	}
	result := ScoreMock(card)

	if result.TotalScore != 0.0 {
		t.Errorf("expected total score 0.0, got %.2f", result.TotalScore)
	}
	if result.HireSignal != "no-hire" {
		t.Errorf("expected no-hire, got %s", result.HireSignal)
	}
	if len(result.WeakAreas) != 7 {
		t.Errorf("expected 7 weak areas for zero scores, got %d: %v", len(result.WeakAreas), result.WeakAreas)
	}
}

func TestScoreMockHireSignalBoundaries(t *testing.T) {
	// Score of exactly 85 should be strong-hire.
	// Score of exactly 65 should be hire.
	// Score of exactly 45 should be mixed.
	// Score of exactly 44 should be no-hire.

	cases := []struct {
		name      string
		card      MockScorecard
		wantRange string // expected hire signal
	}{
		{
			name:      "strong-hire boundary at 85",
			card:      MockScorecard{Clarification: 85, Capacity: 85, Architecture: 85, DeepDive: 85, FailureModes: 85, Tradeoffs: 85, Observability: 85},
			wantRange: "strong-hire",
		},
		{
			name:      "hire range",
			card:      MockScorecard{Clarification: 70, Capacity: 70, Architecture: 70, DeepDive: 70, FailureModes: 70, Tradeoffs: 70, Observability: 70},
			wantRange: "hire",
		},
		{
			name:      "mixed range",
			card:      MockScorecard{Clarification: 50, Capacity: 50, Architecture: 50, DeepDive: 50, FailureModes: 50, Tradeoffs: 50, Observability: 50},
			wantRange: "mixed",
		},
		{
			name:      "no-hire range",
			card:      MockScorecard{Clarification: 20, Capacity: 20, Architecture: 20, DeepDive: 20, FailureModes: 20, Tradeoffs: 20, Observability: 20},
			wantRange: "no-hire",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ScoreMock(tc.card)
			if result.HireSignal != tc.wantRange {
				t.Errorf("card %v: expected %s, got %s (total=%.2f)", tc.card, tc.wantRange, result.HireSignal, result.TotalScore)
			}
		})
	}
}

func TestScoreMockWeakAreaDetection(t *testing.T) {
	card := MockScorecard{
		Clarification: 90,
		Capacity:      90,
		Architecture:  90,
		DeepDive:      90,
		FailureModes:  50, // weak
		Tradeoffs:     40, // weak
		Observability: 30, // weak
	}
	result := ScoreMock(card)

	if len(result.WeakAreas) != 3 {
		t.Errorf("expected 3 weak areas, got %d: %v", len(result.WeakAreas), result.WeakAreas)
	}

	// Verify the specific weak dimensions are flagged.
	weakSet := map[string]bool{}
	for _, w := range result.WeakAreas {
		weakSet[w] = true
	}
	for _, expected := range []string{"failure_modes", "tradeoffs", "observability"} {
		if !weakSet[expected] {
			t.Errorf("expected %s to be in weak areas, got %v", expected, result.WeakAreas)
		}
	}
}

func TestScoreMockBreakdownSumsToTotal(t *testing.T) {
	card := MockScorecard{
		Clarification: 80,
		Capacity:      60,
		Architecture:  70,
		DeepDive:      85,
		FailureModes:  75,
		Tradeoffs:     65,
		Observability: 90,
	}
	result := ScoreMock(card)

	sum := 0.0
	for _, v := range result.Breakdown {
		sum += v
	}

	const epsilon = 1e-9
	if abs64(sum-result.TotalScore) > epsilon {
		t.Errorf("breakdown sum %.4f does not match total score %.4f", sum, result.TotalScore)
	}
}

func TestScoreMockClampAbove100(t *testing.T) {
	// Scores above 100 should be clamped to 100.
	card := MockScorecard{
		Clarification: 200,
		Capacity:      150,
		Architecture:  110,
		DeepDive:      999,
		FailureModes:  100,
		Tradeoffs:     100,
		Observability: 100,
	}
	result := ScoreMock(card)

	// Maximum possible is 100.0.
	if result.TotalScore > 100.0 {
		t.Errorf("total score should not exceed 100 after clamping, got %.2f", result.TotalScore)
	}
}

func TestScoreMockBreakdownKeysPresent(t *testing.T) {
	card := MockScorecard{
		Clarification: 75,
		Capacity:      75,
		Architecture:  75,
		DeepDive:      75,
		FailureModes:  75,
		Tradeoffs:     75,
		Observability: 75,
	}
	result := ScoreMock(card)

	expectedKeys := []string{
		"clarification", "capacity", "architecture",
		"deep_dive", "failure_modes", "tradeoffs", "observability",
	}
	for _, key := range expectedKeys {
		if _, ok := result.Breakdown[key]; !ok {
			t.Errorf("expected breakdown key %q not found in %v", key, result.Breakdown)
		}
	}
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
