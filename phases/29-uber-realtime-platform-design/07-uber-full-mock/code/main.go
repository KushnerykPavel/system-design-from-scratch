package main

import (
	"fmt"
	"strings"
)

// UberScorecard holds raw scores (0–100) for each rubric dimension.
// Weights: GeospatialDesign 20, LocationPipeline 20, DispatchAlgorithm 15,
//          TripStateManagement 15, FailureRecovery 15, Observability 15.
type UberScorecard struct {
	GeospatialDesign    int // 0–100
	LocationPipeline    int // 0–100
	DispatchAlgorithm   int // 0–100
	TripStateManagement int // 0–100
	FailureRecovery     int // 0–100
	Observability       int // 0–100
}

// dimensionWeight pairs a human-readable name with its rubric weight and score.
type dimensionWeight struct {
	Name   string
	Weight int
	Score  int
}

// dimensions returns the ordered list of rubric dimensions with their weights and scores.
func (c UberScorecard) dimensions() []dimensionWeight {
	return []dimensionWeight{
		{"Geospatial Design", 20, c.GeospatialDesign},
		{"Location Pipeline", 20, c.LocationPipeline},
		{"Dispatch Algorithm", 15, c.DispatchAlgorithm},
		{"Trip State Management", 15, c.TripStateManagement},
		{"Failure Recovery", 15, c.FailureRecovery},
		{"Observability", 15, c.Observability},
	}
}

// MockResult holds the output of a mock interview evaluation.
type MockResult struct {
	WeightedScore float64
	HireSignal    string
	WeakAreas     []string
}

// EvaluateUberMock computes a weighted score across all rubric dimensions,
// assigns a hire signal, and identifies weak areas (score < 70 on any dimension).
func EvaluateUberMock(card UberScorecard) MockResult {
	dims := card.dimensions()

	var weightedSum float64
	var totalWeight int
	var weakAreas []string

	for _, d := range dims {
		// Clamp score to 0–100.
		score := d.Score
		if score < 0 {
			score = 0
		}
		if score > 100 {
			score = 100
		}
		weightedSum += float64(score) * float64(d.Weight)
		totalWeight += d.Weight

		if score < 70 {
			weakAreas = append(weakAreas, d.Name)
		}
	}

	weighted := weightedSum / float64(totalWeight)

	signal := hireSignal(weighted)

	return MockResult{
		WeightedScore: weighted,
		HireSignal:    signal,
		WeakAreas:     weakAreas,
	}
}

// hireSignal maps a weighted percentage score to a hire band.
func hireSignal(score float64) string {
	switch {
	case score >= 85:
		return "Strong Hire"
	case score >= 70:
		return "Hire"
	case score >= 55:
		return "Weak Hire"
	default:
		return "No Hire"
	}
}

func main() {
	candidates := []struct {
		name string
		card UberScorecard
	}{
		{
			name: "Alice (strong hire)",
			card: UberScorecard{
				GeospatialDesign:    95,
				LocationPipeline:    90,
				DispatchAlgorithm:   85,
				TripStateManagement: 88,
				FailureRecovery:     82,
				Observability:       80,
			},
		},
		{
			name: "Bob (hire — weak on observability)",
			card: UberScorecard{
				GeospatialDesign:    80,
				LocationPipeline:    75,
				DispatchAlgorithm:   72,
				TripStateManagement: 78,
				FailureRecovery:     70,
				Observability:       45,
			},
		},
		{
			name: "Carol (weak hire — no state machine or failure modes)",
			card: UberScorecard{
				GeospatialDesign:    70,
				LocationPipeline:    65,
				DispatchAlgorithm:   60,
				TripStateManagement: 40,
				FailureRecovery:     35,
				Observability:       50,
			},
		},
		{
			name: "Dan (no hire — database for driver locations)",
			card: UberScorecard{
				GeospatialDesign:    20,
				LocationPipeline:    25,
				DispatchAlgorithm:   30,
				TripStateManagement: 20,
				FailureRecovery:     15,
				Observability:       10,
			},
		},
	}

	for _, c := range candidates {
		result := EvaluateUberMock(c.card)
		fmt.Printf("Candidate: %s\n", c.name)
		fmt.Printf("  Weighted Score: %.1f%%\n", result.WeightedScore)
		fmt.Printf("  Signal:         %s\n", result.HireSignal)
		if len(result.WeakAreas) > 0 {
			fmt.Printf("  Weak Areas:     %s\n", strings.Join(result.WeakAreas, ", "))
		} else {
			fmt.Printf("  Weak Areas:     none\n")
		}
		fmt.Println()
	}
}
