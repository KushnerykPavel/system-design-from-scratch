package main

import (
	"encoding/json"
	"os"
)

type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{Category: "working-backwards", Description: "started from customer experience, not technical solution", Observed: false},
		{Category: "lp-demonstration", Description: "demonstrated at least one Amazon Leadership Principle in framing", Observed: false},
		{Category: "service-decomposition", Description: "decomposed into independently deployable two-pizza-team services", Observed: false},
		{Category: "capacity", Description: "gave concrete capacity numbers before architecture", Observed: false},
		{Category: "aws-services", Description: "selected appropriate AWS services with justification", Observed: false},
		{Category: "failure-modes", Description: "identified top failure modes with mitigations", Observed: false},
		{Category: "observability", Description: "specified concrete metrics and alarms", Observed: false},
		{Category: "tradeoffs", Description: "named trade-offs with explicitly rejected alternatives", Observed: false},
	}
}

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
	for i := range signals {
		if signals[i].Category != "lp-demonstration" {
			signals[i].Observed = true
		}
	}
	result := EvaluateRubric(signals)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
