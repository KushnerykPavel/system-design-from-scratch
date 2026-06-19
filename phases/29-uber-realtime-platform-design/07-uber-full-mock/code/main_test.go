package main

import (
	"testing"
)

// TestStrongHireScore verifies that near-perfect scores produce "Strong Hire".
func TestStrongHireScore(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    95,
		LocationPipeline:    90,
		DispatchAlgorithm:   90,
		TripStateManagement: 90,
		FailureRecovery:     85,
		Observability:       85,
	}
	result := EvaluateUberMock(card)
	if result.HireSignal != "Strong Hire" {
		t.Errorf("expected Strong Hire, got %s (score=%.1f)", result.HireSignal, result.WeightedScore)
	}
	if len(result.WeakAreas) != 0 {
		t.Errorf("expected no weak areas, got %v", result.WeakAreas)
	}
}

// TestHireScore verifies mid-range scores produce "Hire".
func TestHireScore(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    80,
		LocationPipeline:    75,
		DispatchAlgorithm:   72,
		TripStateManagement: 75,
		FailureRecovery:     70,
		Observability:       70,
	}
	result := EvaluateUberMock(card)
	if result.HireSignal != "Hire" {
		t.Errorf("expected Hire, got %s (score=%.1f)", result.HireSignal, result.WeightedScore)
	}
}

// TestWeakHireScore verifies that scores in the 55–69 band produce "Weak Hire".
func TestWeakHireScore(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    65,
		LocationPipeline:    60,
		DispatchAlgorithm:   58,
		TripStateManagement: 55,
		FailureRecovery:     50,
		Observability:       55,
	}
	result := EvaluateUberMock(card)
	if result.HireSignal != "Weak Hire" {
		t.Errorf("expected Weak Hire, got %s (score=%.1f)", result.HireSignal, result.WeightedScore)
	}
}

// TestNoHireScore verifies very low scores produce "No Hire".
func TestNoHireScore(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    20,
		LocationPipeline:    25,
		DispatchAlgorithm:   30,
		TripStateManagement: 20,
		FailureRecovery:     15,
		Observability:       10,
	}
	result := EvaluateUberMock(card)
	if result.HireSignal != "No Hire" {
		t.Errorf("expected No Hire, got %s (score=%.1f)", result.HireSignal, result.WeightedScore)
	}
}

// TestWeakAreasIdentified verifies that dimensions below 70 are flagged.
func TestWeakAreasIdentified(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    90,
		LocationPipeline:    85,
		DispatchAlgorithm:   80,
		TripStateManagement: 40, // weak
		FailureRecovery:     35, // weak
		Observability:       75,
	}
	result := EvaluateUberMock(card)

	weakSet := make(map[string]bool)
	for _, w := range result.WeakAreas {
		weakSet[w] = true
	}

	if !weakSet["Trip State Management"] {
		t.Error("expected Trip State Management in weak areas")
	}
	if !weakSet["Failure Recovery"] {
		t.Error("expected Failure Recovery in weak areas")
	}
	if weakSet["Geospatial Design"] {
		t.Error("Geospatial Design (90) should not be a weak area")
	}
}

// TestWeightedScoreAccuracy verifies that weights are applied correctly.
// All scores = 100 should produce weighted score = 100.0.
func TestWeightedScoreAllPerfect(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    100,
		LocationPipeline:    100,
		DispatchAlgorithm:   100,
		TripStateManagement: 100,
		FailureRecovery:     100,
		Observability:       100,
	}
	result := EvaluateUberMock(card)
	if result.WeightedScore != 100.0 {
		t.Errorf("expected 100.0 weighted score, got %.2f", result.WeightedScore)
	}
}

// TestWeightedScoreAllZero verifies that all-zero scores produce 0.0 and "No Hire".
func TestWeightedScoreAllZero(t *testing.T) {
	card := UberScorecard{}
	result := EvaluateUberMock(card)
	if result.WeightedScore != 0.0 {
		t.Errorf("expected 0.0 weighted score, got %.2f", result.WeightedScore)
	}
	if result.HireSignal != "No Hire" {
		t.Errorf("expected No Hire for zero scores, got %s", result.HireSignal)
	}
}

// TestScoreClampedAbove100 verifies that out-of-range scores are clamped to 100.
func TestScoreClampedAbove100(t *testing.T) {
	card := UberScorecard{
		GeospatialDesign:    150, // over-range
		LocationPipeline:    100,
		DispatchAlgorithm:   100,
		TripStateManagement: 100,
		FailureRecovery:     100,
		Observability:       100,
	}
	result := EvaluateUberMock(card)
	if result.WeightedScore != 100.0 {
		t.Errorf("expected clamped score 100.0, got %.2f", result.WeightedScore)
	}
}

// TestBoundaryHireAt70 verifies that score exactly 70.0 maps to "Hire".
func TestBoundaryHireAt70(t *testing.T) {
	// All dimensions at 70 → weighted score = 70.
	card := UberScorecard{
		GeospatialDesign:    70,
		LocationPipeline:    70,
		DispatchAlgorithm:   70,
		TripStateManagement: 70,
		FailureRecovery:     70,
		Observability:       70,
	}
	result := EvaluateUberMock(card)
	if result.HireSignal != "Hire" {
		t.Errorf("expected Hire at score 70.0, got %s", result.HireSignal)
	}
}
