package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// DesignSignal represents a candidate response signal in a Meta interview rubric.
type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

// RubricResult summarizes how a candidate scored against the Meta rubric.
type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

// DefaultSignals returns the standard Meta strong-hire signals.
func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{Category: "social-graph-scale", Description: "named graph scale (3B nodes, 100B+ edges) and its structural implications", Observed: false},
		{Category: "tao-cache-awareness", Description: "referenced TAO or two-level graph cache with objects and associations", Observed: false},
		{Category: "fanout-strategy", Description: "distinguished push-on-write vs pull-on-read with a follower-count threshold", Observed: false},
		{Category: "privacy-compliance", Description: "identified privacy enforcement before fan-out, not eventually", Observed: false},
		{Category: "ml-integration", Description: "described a ranking or recommendation integration point with feature store", Observed: false},
		{Category: "failure-first", Description: "named top failure modes and explicit mitigations", Observed: false},
		{Category: "observability", Description: "specified at least two concrete metrics or traces", Observed: false},
		{Category: "meta-vocab", Description: "used at least one Meta-specific term correctly (TAO, Haystack, Scuba, Iris, Everstore)", Observed: false},
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
	// Simulate a strong-hire candidate: all signals observed except meta-vocab.
	for i := range signals {
		if signals[i].Category != "meta-vocab" {
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
