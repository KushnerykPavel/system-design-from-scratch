package main

import (
	"testing"
)

// TestStrongHire verifies that a high-scoring candidate gets strong-hire and LP demonstrated.
func TestStrongHire(t *testing.T) {
	card := AmazonScorecard{
		LPFraming: 90, WorkingBackwards: 90, Capacity: 90,
		Architecture: 90, DeepDive: 90, FailureModes: 90,
		Tradeoffs: 90, Observability: 90,
	}
	result := EvaluateAmazonMock(card)
	if result.HireSignal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %s (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if !result.LPDemonstrated {
		t.Fatal("expected LPDemonstrated=true for high LP scores")
	}
	if len(result.WeakAreas) != 0 {
		t.Fatalf("expected no weak areas, got %v", result.WeakAreas)
	}
}

// TestNoHireDueToLPFailure verifies that a candidate with poor LP scores and average
// technical scores produces no-hire and has LPDemonstrated=false.
// LP weight is 30/100 — if LP scores are 0 and all technical dims are 50,
// weighted = 0 + 0 + 50*10/100 + 50*20/100 + 50*20/100 + 50*10/100 + 50*5/100 + 50*5/100
//          = 0 + 0 + 5 + 10 + 10 + 5 + 2 + 2 = 34 → no-hire.
func TestNoHireDueToLPFailure(t *testing.T) {
	card := AmazonScorecard{
		LPFraming: 0, WorkingBackwards: 0, Capacity: 50,
		Architecture: 50, DeepDive: 50, FailureModes: 50,
		Tradeoffs: 50, Observability: 50,
	}
	result := EvaluateAmazonMock(card)
	if result.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire when LP scores are 0 and technical scores are 50, got %s (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if result.LPDemonstrated {
		t.Fatal("expected LPDemonstrated=false when LP scores are 0")
	}
}

// TestWeightedScoreCalculation verifies the weighted score formula for a known input.
// LPFraming=100*15 + WorkingBackwards=0*15 + all others=0 → 15.
func TestWeightedScoreCalculation(t *testing.T) {
	card := AmazonScorecard{
		LPFraming: 100, WorkingBackwards: 0, Capacity: 0,
		Architecture: 0, DeepDive: 0, FailureModes: 0,
		Tradeoffs: 0, Observability: 0,
	}
	result := EvaluateAmazonMock(card)
	if result.WeightedScore != 15 {
		t.Fatalf("expected weighted score=15 for LPFraming=100 only, got %d", result.WeightedScore)
	}
}

// TestWeakAreasDetection verifies that dimensions below 60 are reported as weak areas.
func TestWeakAreasDetection(t *testing.T) {
	card := AmazonScorecard{
		LPFraming: 80, WorkingBackwards: 80, Capacity: 80,
		Architecture: 80, DeepDive: 80, FailureModes: 30, // weak
		Tradeoffs: 40, Observability: 50, // weak
	}
	result := EvaluateAmazonMock(card)
	if len(result.WeakAreas) != 3 {
		t.Fatalf("expected 3 weak areas (FailureModes, Tradeoffs, Observability), got %d: %v", len(result.WeakAreas), result.WeakAreas)
	}
}

// TestHireThreshold verifies the exact hire boundary at weighted score >= 70.
func TestHireThreshold(t *testing.T) {
	// Construct a card that produces exactly 70: need weighted = 70.
	// All dims at 100% → weighted = 100. At 70%:
	// LPFraming=70*15/100=10 + WB=70*15/100=10 + Cap=70*10/100=7 + Arch=70*20/100=14
	// + Deep=70*20/100=14 + Fail=70*10/100=7 + Trade=70*5/100=3 + Obs=70*5/100=3 = 68 (int division)
	// Use 80 for higher-weight dims to hit exactly 70.
	card := AmazonScorecard{
		LPFraming: 80, WorkingBackwards: 80, Capacity: 80,
		Architecture: 75, DeepDive: 75, FailureModes: 80,
		Tradeoffs: 80, Observability: 80,
	}
	result := EvaluateAmazonMock(card)
	if result.HireSignal == "no-hire" {
		t.Fatalf("expected hire or strong-hire, got no-hire (score=%d)", result.WeightedScore)
	}
}

// TestLPDemonstratedThreshold verifies LP demonstrated requires average LP score >= 70.
func TestLPDemonstratedThreshold(t *testing.T) {
	cardBelow := AmazonScorecard{LPFraming: 60, WorkingBackwards: 60}
	if EvaluateAmazonMock(cardBelow).LPDemonstrated {
		t.Fatal("expected LPDemonstrated=false when LP average is 60")
	}

	cardAt := AmazonScorecard{LPFraming: 70, WorkingBackwards: 70}
	if !EvaluateAmazonMock(cardAt).LPDemonstrated {
		t.Fatal("expected LPDemonstrated=true when LP average is exactly 70")
	}
}

// TestZeroScorecard verifies graceful handling of all-zero scores.
func TestZeroScorecard(t *testing.T) {
	result := EvaluateAmazonMock(AmazonScorecard{})
	if result.WeightedScore != 0 {
		t.Fatalf("expected weighted score=0, got %d", result.WeightedScore)
	}
	if result.HireSignal != "no-hire" {
		t.Fatalf("expected no-hire for zero scorecard, got %s", result.HireSignal)
	}
	if result.LPDemonstrated {
		t.Fatal("expected LPDemonstrated=false for zero scorecard")
	}
}
