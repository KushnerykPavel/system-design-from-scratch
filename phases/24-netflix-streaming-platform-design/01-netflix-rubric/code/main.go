package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// DesignSignal represents a candidate response signal in a Netflix interview rubric.
type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

// RubricResult summarizes how a candidate scored against the Netflix rubric.
type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

// DefaultSignals returns the standard Netflix strong-hire signals.
func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{Category: "capacity", Description: "gave capacity numbers before architecture", Observed: false},
		{Category: "decomposition", Description: "named specific services with ownership boundaries", Observed: false},
		{Category: "failure", Description: "identified top failure modes explicitly", Observed: false},
		{Category: "fallback", Description: "described fallback chain for at least one dependency", Observed: false},
		{Category: "experimentation", Description: "mentioned measurement or A/B testing hook", Observed: false},
		{Category: "tradeoffs", Description: "named trade-offs with rejected alternatives", Observed: false},
		{Category: "observability", Description: "specified at least two concrete metrics or traces", Observed: false},
		{Category: "netflix-vocab", Description: "used at least one Netflix-specific term correctly", Observed: false},
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
	// Simulate a strong-hire candidate: all signals present except netflix-vocab.
	for i := range signals {
		if signals[i].Category != "netflix-vocab" {
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
