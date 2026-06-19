package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// AmazonScorecard holds scores for each rubric dimension (0–100 each).
// Weights: LPFraming 15, WorkingBackwards 15, Capacity 10, Architecture 20,
// DeepDive 20, FailureModes 10, Tradeoffs 5, Observability 5. Total: 100 points.
type AmazonScorecard struct {
	LPFraming       int // 0-100: did candidate demonstrate LPs with specific design decisions?
	WorkingBackwards int // 0-100: was customer promise stated before architecture?
	Capacity        int // 0-100: correct order-of-magnitude numbers derived from first principles?
	Architecture    int // 0-100: three-layer stack, correct data stores, serving path latency?
	DeepDive        int // 0-100: cold start, latency budget, model update cadence?
	FailureModes    int // 0-100: at least 3 failure modes with concrete mitigations?
	Tradeoffs       int // 0-100: at least 2 explicit trade-offs with alternatives named?
	Observability   int // 0-100: SLOs stated, named metrics, business vs technical metric?
}

// dimensionWeights defines the contribution of each dimension to the total score.
// All weights must sum to 100.
var dimensionWeights = map[string]int{
	"LPFraming":        15,
	"WorkingBackwards": 15,
	"Capacity":         10,
	"Architecture":     20,
	"DeepDive":         20,
	"FailureModes":     10,
	"Tradeoffs":        5,
	"Observability":    5,
}

// MockResult holds the evaluation output.
type MockResult struct {
	WeightedScore   int      `json:"weighted_score"`   // 0-100 weighted total
	HireSignal      string   `json:"hire_signal"`      // "strong-hire" / "hire" / "no-hire"
	LPDemonstrated  bool     `json:"lp_demonstrated"`  // true if LPFraming + WorkingBackwards average >= 70
	WeakAreas       []string `json:"weak_areas"`       // dimensions scoring below 60 (raw)
}

// EvaluateAmazonMock computes the weighted score and hire signal from a scorecard.
//
// Hire signal thresholds (weighted score):
//   - strong-hire: >= 85
//   - hire:        >= 70
//   - no-hire:     < 70
//
// LP demonstrated: true when the average of LPFraming and WorkingBackwards >= 70.
// Weak areas: any dimension with a raw score below 60.
func EvaluateAmazonMock(card AmazonScorecard) MockResult {
	scores := map[string]int{
		"LPFraming":        card.LPFraming,
		"WorkingBackwards": card.WorkingBackwards,
		"Capacity":         card.Capacity,
		"Architecture":     card.Architecture,
		"DeepDive":         card.DeepDive,
		"FailureModes":     card.FailureModes,
		"Tradeoffs":        card.Tradeoffs,
		"Observability":    card.Observability,
	}

	// Compute weighted score: sum(score * weight / 100) for each dimension.
	weighted := 0
	for dim, raw := range scores {
		weight := dimensionWeights[dim]
		weighted += (raw * weight) / 100
	}

	// Hire signal.
	hireSignal := "no-hire"
	if weighted >= 85 {
		hireSignal = "strong-hire"
	} else if weighted >= 70 {
		hireSignal = "hire"
	}

	// LP demonstrated: average of LPFraming and WorkingBackwards >= 70.
	lpAvg := (card.LPFraming + card.WorkingBackwards) / 2
	lpDemonstrated := lpAvg >= 70

	// Weak areas: raw score below 60.
	var weakAreas []string
	orderedDims := []string{
		"LPFraming", "WorkingBackwards", "Capacity", "Architecture",
		"DeepDive", "FailureModes", "Tradeoffs", "Observability",
	}
	for _, dim := range orderedDims {
		if scores[dim] < 60 {
			weakAreas = append(weakAreas, fmt.Sprintf("%s (%d/100)", dim, scores[dim]))
		}
	}

	return MockResult{
		WeightedScore:  weighted,
		HireSignal:     hireSignal,
		LPDemonstrated: lpDemonstrated,
		WeakAreas:      weakAreas,
	}
}

func main() {
	examples := []struct {
		label string
		card  AmazonScorecard
	}{
		{
			label: "Strong hire — all dimensions excellent",
			card: AmazonScorecard{
				LPFraming: 90, WorkingBackwards: 90, Capacity: 85,
				Architecture: 90, DeepDive: 90, FailureModes: 85,
				Tradeoffs: 90, Observability: 85,
			},
		},
		{
			label: "Technical depth but no LP alignment — no hire",
			card: AmazonScorecard{
				LPFraming: 20, WorkingBackwards: 20, Capacity: 90,
				Architecture: 90, DeepDive: 90, FailureModes: 80,
				Tradeoffs: 80, Observability: 70,
			},
		},
		{
			label: "Borderline hire — solid architecture, weak failure modes",
			card: AmazonScorecard{
				LPFraming: 75, WorkingBackwards: 75, Capacity: 70,
				Architecture: 75, DeepDive: 70, FailureModes: 40,
				Tradeoffs: 65, Observability: 50,
			},
		},
		{
			label: "Strong hire threshold exact (score=85)",
			card: AmazonScorecard{
				LPFraming: 100, WorkingBackwards: 100, Capacity: 100,
				Architecture: 75, DeepDive: 75, FailureModes: 100,
				Tradeoffs: 100, Observability: 100,
			},
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for _, ex := range examples {
		result := EvaluateAmazonMock(ex.card)
		fmt.Printf("=== %s ===\n", ex.label)
		_ = enc.Encode(result)
		fmt.Println()
	}
}
