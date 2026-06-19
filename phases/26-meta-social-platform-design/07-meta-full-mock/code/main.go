package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Dimension weight constants (must sum to 100).
const (
	weightClarification = 10
	weightCapacity      = 10
	weightArchitecture  = 20
	weightDeepDive      = 20
	weightFailureModes  = 15
	weightTradeoffs     = 15
	weightObservability = 10
)

// MockScorecard holds raw scores (0–100) for each interview dimension.
// Each score represents how well the candidate performed on that dimension,
// independently of its weighting.
type MockScorecard struct {
	Clarification int `json:"clarification"`  // Did they ask the right questions?
	Capacity      int `json:"capacity"`       // Did they estimate numbers?
	Architecture  int `json:"architecture"`   // Did they cover all layers correctly?
	DeepDive      int `json:"deep_dive"`      // Did they go deep on a subsystem?
	FailureModes  int `json:"failure_modes"`  // Did they identify failures + mitigations?
	Tradeoffs     int `json:"tradeoffs"`      // Did they reason about trade-offs?
	Observability int `json:"observability"`  // Did they name metrics and alerts?
}

// InterviewResult is the output of scoring a MockScorecard.
type InterviewResult struct {
	TotalScore  float64  `json:"total_score"`   // 0.0–100.0
	HireSignal  string   `json:"hire_signal"`   // "strong-hire" | "hire" | "mixed" | "no-hire"
	WeakAreas   []string `json:"weak_areas"`    // dimensions where score < 60
	Breakdown   map[string]float64 `json:"breakdown"` // weighted contribution per dimension
}

// clamp ensures a score stays within [0, 100].
func clamp(score int) float64 {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return float64(score)
}

// weightedContribution returns the weighted score for a dimension.
// raw is 0–100; weight is the dimension's point allocation out of 100.
func weightedContribution(raw int, weight int) float64 {
	return clamp(raw) * float64(weight) / 100.0
}

// ScoreMock computes the final interview result from a scorecard.
func ScoreMock(card MockScorecard) InterviewResult {
	breakdown := map[string]float64{
		"clarification": weightedContribution(card.Clarification, weightClarification),
		"capacity":      weightedContribution(card.Capacity, weightCapacity),
		"architecture":  weightedContribution(card.Architecture, weightArchitecture),
		"deep_dive":     weightedContribution(card.DeepDive, weightDeepDive),
		"failure_modes": weightedContribution(card.FailureModes, weightFailureModes),
		"tradeoffs":     weightedContribution(card.Tradeoffs, weightTradeoffs),
		"observability": weightedContribution(card.Observability, weightObservability),
	}

	total := 0.0
	for _, v := range breakdown {
		total += v
	}

	signal := "no-hire"
	switch {
	case total >= 85:
		signal = "strong-hire"
	case total >= 65:
		signal = "hire"
	case total >= 45:
		signal = "mixed"
	}

	// Identify weak areas: dimensions where the raw score is below 60.
	weakAreas := []string{}
	type dimScore struct {
		name  string
		score int
	}
	dims := []dimScore{
		{"clarification", card.Clarification},
		{"capacity", card.Capacity},
		{"architecture", card.Architecture},
		{"deep_dive", card.DeepDive},
		{"failure_modes", card.FailureModes},
		{"tradeoffs", card.Tradeoffs},
		{"observability", card.Observability},
	}
	for _, d := range dims {
		if d.score < 60 {
			weakAreas = append(weakAreas, d.name)
		}
	}

	return InterviewResult{
		TotalScore: total,
		HireSignal: signal,
		WeakAreas:  weakAreas,
		Breakdown:  breakdown,
	}
}

func main() {
	// Scenario A: Strong candidate — missed some observability depth.
	scenarioA := MockScorecard{
		Clarification: 90, // asked 4 good questions with numbers
		Capacity:      85, // estimated DAU, RPS, storage, fanout ratio
		Architecture:  85, // covered all 4 layers with Meta-specific components
		DeepDive:      90, // went deep on ranking: 3-stage funnel, FAISS, feature store
		FailureModes:  80, // named 3 failures with mitigations
		Tradeoffs:     85, // named push/pull fanout, async transcoding, three-stage funnel trade-offs
		Observability: 55, // named 2 metrics but no alert thresholds
	}

	// Scenario B: Mixed candidate — good breadth, shallow depth, no failure modes.
	scenarioB := MockScorecard{
		Clarification: 70, // asked 2 questions, no numbers
		Capacity:      40, // rough estimates only
		Architecture:  75, // covered 3 of 4 layers
		DeepDive:      45, // stayed high-level on ranking — "we use ML"
		FailureModes:  30, // mentioned retries but no specific failures
		Tradeoffs:     60, // named push vs pull fanout but no rejected alternatives
		Observability: 50, // named latency metric only
	}

	resultA := ScoreMock(scenarioA)
	resultB := ScoreMock(scenarioB)

	output := map[string]any{
		"scenario_a_strong_candidate": map[string]any{
			"scorecard": scenarioA,
			"result":    resultA,
		},
		"scenario_b_mixed_candidate": map[string]any{
			"scorecard": scenarioB,
			"result":    resultB,
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
