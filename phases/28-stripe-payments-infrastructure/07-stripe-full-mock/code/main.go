package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// StripeConnectScorecard holds scores (0-100) for each evaluation dimension.
type StripeConnectScorecard struct {
	Idempotency     int // 20 pts weight
	LedgerDesign    int // 20 pts weight
	APIDesign       int // 15 pts weight
	FraudHooks      int // 15 pts weight
	FailureRecovery int // 15 pts weight
	Observability   int // 15 pts weight
}

// MockResult holds the evaluation output for a mock answer.
type MockResult struct {
	WeightedScore int      `json:"weighted_score"`
	HireSignal    string   `json:"hire_signal"`
	WeakAreas     []string `json:"weak_areas"`
}

// weights maps each dimension to its contribution to the total (out of 100).
var weights = map[string]int{
	"Idempotency":     20,
	"LedgerDesign":    20,
	"APIDesign":       15,
	"FraudHooks":      15,
	"FailureRecovery": 15,
	"Observability":   15,
}

// EvaluateMock scores a StripeConnectScorecard and returns a MockResult.
// WeightedScore = sum of (dimension_score * weight / 100) for each dimension.
// Hire thresholds: >=85 strong-hire, >=65 hire, >=45 mixed, <45 no-hire.
// WeakAreas: any dimension scoring below 50.
func EvaluateMock(card StripeConnectScorecard) MockResult {
	scores := map[string]int{
		"Idempotency":     card.Idempotency,
		"LedgerDesign":    card.LedgerDesign,
		"APIDesign":       card.APIDesign,
		"FraudHooks":      card.FraudHooks,
		"FailureRecovery": card.FailureRecovery,
		"Observability":   card.Observability,
	}

	total := 0
	for dim, score := range scores {
		w := weights[dim]
		total += score * w / 100
	}

	hireSignal := "no-hire"
	switch {
	case total >= 85:
		hireSignal = "strong-hire"
	case total >= 65:
		hireSignal = "hire"
	case total >= 45:
		hireSignal = "mixed"
	}

	var weakAreas []string
	// Order matters for deterministic output; check in fixed order.
	dimOrder := []string{
		"Idempotency", "LedgerDesign", "APIDesign",
		"FraudHooks", "FailureRecovery", "Observability",
	}
	for _, dim := range dimOrder {
		if scores[dim] < 50 {
			weakAreas = append(weakAreas, dim)
		}
	}

	return MockResult{
		WeightedScore: total,
		HireSignal:    hireSignal,
		WeakAreas:     weakAreas,
	}
}

func main() {
	// Sample answer: strong on idempotency and ledger, weak on observability and fraud hooks.
	sample := StripeConnectScorecard{
		Idempotency:     90,
		LedgerDesign:    85,
		APIDesign:       70,
		FraudHooks:      40,
		FailureRecovery: 75,
		Observability:   35,
	}

	result := EvaluateMock(sample)

	fmt.Println("=== Mock Interview Scorecard ===")
	fmt.Printf("Idempotency     : %d/100  (weight: 20%%)\n", sample.Idempotency)
	fmt.Printf("LedgerDesign    : %d/100  (weight: 20%%)\n", sample.LedgerDesign)
	fmt.Printf("APIDesign       : %d/100  (weight: 15%%)\n", sample.APIDesign)
	fmt.Printf("FraudHooks      : %d/100  (weight: 15%%)\n", sample.FraudHooks)
	fmt.Printf("FailureRecovery : %d/100  (weight: 15%%)\n", sample.FailureRecovery)
	fmt.Printf("Observability   : %d/100  (weight: 15%%)\n", sample.Observability)
	fmt.Println()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
