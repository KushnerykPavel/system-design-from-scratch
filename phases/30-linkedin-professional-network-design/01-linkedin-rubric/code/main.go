package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// DesignSignal represents a candidate response signal in a LinkedIn interview rubric.
type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

// RubricResult summarizes how a candidate scored against the LinkedIn rubric.
type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

// DefaultSignals returns the standard LinkedIn strong-hire signals.
func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{Category: "graph-traversal", Description: "named graph scale (950M nodes, 475B+ edges) and BFS traversal cost with precomputation rationale", Observed: false},
		{Category: "kafka-pipeline-awareness", Description: "LinkedIn invented Kafka; described partitioning strategy, consumer groups, and at-least-once semantics", Observed: false},
		{Category: "feed-ranking", Description: "distinguished professional relevance signals from social virality; named multi-stage ranking pipeline", Observed: false},
		{Category: "privacy-gdpr", Description: "identified GDPR right-to-be-forgotten deletion pipeline propagating to Espresso, Kafka, Pinot, and search", Observed: false},
		{Category: "job-matching", Description: "described skills ontology, semantic matching, and recruiter vs member matching asymmetry for job recommendations", Observed: false},
		{Category: "member-data-at-scale", Description: "designed for 950M heterogeneous member profiles across Espresso, Venice feature store, and search index", Observed: false},
		{Category: "observability", Description: "named at least two concrete metrics: Kafka consumer lag and PYMK latency p99", Observed: false},
		{Category: "linkedin-vocab", Description: "used at least one LinkedIn-specific term correctly (Espresso, Venice, Pinot, Samza, Azkaban, Brooklin, PYMK)", Observed: false},
	}
}

// EvaluateRubric scores the observed signals and returns a result.
func EvaluateRubric(signals []DesignSignal) RubricResult {
	total := len(signals)
	observed := 0
	var missing []string
	for _, s := range signals {
		if s.Observed {
			observed++
		} else {
			missing = append(missing, s.Category+": "+s.Description)
		}
	}
	score := 0
	if total > 0 {
		score = (observed * 100) / total
	}
	hire := "no-hire"
	switch {
	case score >= 88:
		hire = "strong-hire"
	case score >= 63:
		hire = "hire"
	case score >= 38:
		hire = "mixed"
	}
	return RubricResult{
		TotalSignals:    total,
		ObservedSignals: observed,
		ScorePercent:    score,
		HireSignal:      hire,
		MissingAreas:    missing,
	}
}

func main() {
	signals := DefaultSignals()
	// Simulate a candidate who missed the privacy-gdpr signal.
	for i := range signals {
		if signals[i].Category != "privacy-gdpr" {
			signals[i].Observed = true
		}
	}
	result := EvaluateRubric(signals)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
