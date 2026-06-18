package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Dimension names used in the rubric.
const (
	DimClarification   = "clarification"
	DimBreadth         = "breadth"
	DimDepth           = "depth"
	DimFailureModes    = "failure_modes"
	DimPivotHandling   = "pivot_handling"
	DimTradeoffs       = "tradeoffs"
	DimNetflixVocab    = "netflix_vocabulary"
	DimSummary         = "summary"
)

// Score holds the rating for one rubric dimension.
type Score struct {
	Dimension string `json:"dimension"`
	Points    int    `json:"points"` // 0, 1, or 2
	Note      string `json:"note"`
}

// MockResult is the final scored rubric output.
type MockResult struct {
	TotalPoints int     `json:"total_points"`
	MaxPoints   int     `json:"max_points"`
	Signal      string  `json:"hire_signal"`
	Scores      []Score `json:"scores"`
}

// EvaluateMock computes the rubric result from a set of dimension scores.
func EvaluateMock(scores []Score) MockResult {
	total := 0
	for _, s := range scores {
		total += s.Points
	}
	maxPoints := len(scores) * 2
	signal := "no-hire"
	switch {
	case total >= 14:
		signal = "strong-hire"
	case total >= 10:
		signal = "hire"
	case total >= 6:
		signal = "mixed"
	}
	return MockResult{
		TotalPoints: total,
		MaxPoints:   maxPoints,
		Signal:      signal,
		Scores:      scores,
	}
}

func main() {
	// Example: simulate a candidate who scored well overall but missed the pivot.
	scores := []Score{
		{Dimension: DimClarification, Points: 2, Note: "Asked 4 targeted questions, provided subscriber and CDN capacity numbers"},
		{Dimension: DimBreadth, Points: 2, Note: "Named all 4 layers: encoding, CDN, ABR, recommendation"},
		{Dimension: DimDepth, Points: 2, Note: "Went deep on OCA cache hierarchy with failure modes and eviction policy"},
		{Dimension: DimFailureModes, Points: 2, Note: "Named CDN miss, recommendation fallback, region failure, and encoding crash"},
		{Dimension: DimPivotHandling, Points: 1, Note: "Answered live streaming pivot but lost track of overall design structure"},
		{Dimension: DimTradeoffs, Points: 2, Note: "Named per-title encoding vs fixed ladder, at-least-once vs exactly-once"},
		{Dimension: DimNetflixVocab, Points: 2, Note: "Used Open Connect, EVCache, Chaos Monkey, FIT correctly"},
		{Dimension: DimSummary, Points: 1, Note: "Partial summary: mentioned top trade-offs but skipped observability next steps"},
	}

	result := EvaluateMock(scores)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
