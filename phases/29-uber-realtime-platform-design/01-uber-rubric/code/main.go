package main

import (
	"encoding/json"
	"os"
)

// DesignSignal represents one observable behaviour during a system design interview.
type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

// RubricResult summarises the evaluation of a candidate's answer.
type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

// DefaultSignals returns the eight Uber-specific design signals that separate
// strong-hire from generic real-time system answers.
func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{
			Category:    "geospatial-indexing",
			Description: "named H3 hexagonal cells and k-ring proximity query instead of PostGIS or Redis GEO",
			Observed:    false,
		},
		{
			Category:    "realtime-location-pipeline",
			Description: "described Kafka ingest at 20M events/sec partitioned by driver ID feeding Redis cell store",
			Observed:    false,
		},
		{
			Category:    "dispatch-decomposition",
			Description: "distinguished DISCO batch matching (500ms window, Hungarian-variant) from greedy nearest",
			Observed:    false,
		},
		{
			Category:    "surge-pricing-design",
			Description: "designed per-H3-cell demand ratio with EMA dampening and recalculation cadence",
			Observed:    false,
		},
		{
			Category:    "trip-state-machine",
			Description: "drew driver and trip state machines with all transitions and recovery paths",
			Observed:    false,
		},
		{
			Category:    "failure-recovery",
			Description: "addressed Redis shard failure, Kafka consumer lag, and driver offline mid-trip",
			Observed:    false,
		},
		{
			Category:    "observability",
			Description: "named P99 ETA error, match rate, driver utilization, and Kafka lag as key metrics",
			Observed:    false,
		},
		{
			Category:    "uber-vocab",
			Description: "used at least two: H3, DISCO, Cadence/Temporal, Ringpop, TChannel, M3, Docstore",
			Observed:    false,
		},
	}
}

// EvaluateRubric scores a candidate based on which signals were observed.
// Thresholds: >=88% strong-hire, >=63% hire, >=38% mixed, <38% no-hire.
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
	// Simulate a strong candidate who covered everything except surge pricing design.
	signals := DefaultSignals()
	for i := range signals {
		if signals[i].Category != "surge-pricing-design" {
			signals[i].Observed = true
		}
	}
	result := EvaluateRubric(signals)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
