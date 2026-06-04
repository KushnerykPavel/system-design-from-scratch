package main

import (
	"encoding/json"
	"flag"
	"os"
)

type DrillScorecard struct {
	Name                     string `json:"name"`
	DefinesInvariant         bool   `json:"defines_invariant"`
	NamesSourceOfTruth       bool   `json:"names_source_of_truth"`
	SizesRetryAmplification  bool   `json:"sizes_retry_amplification"`
	HasDeepDive              bool   `json:"has_deep_dive"`
	CoversFailureModes       bool   `json:"covers_failure_modes"`
	CoversObservability      bool   `json:"covers_observability"`
	HandlesRedesign          bool   `json:"handles_redesign"`
}

func ValidateDrillScorecard(card DrillScorecard) []string {
	var issues []string
	if !card.DefinesInvariant {
		issues = append(issues, "defines_invariant should be true so the answer starts with a real correctness boundary")
	}
	if !card.NamesSourceOfTruth {
		issues = append(issues, "names_source_of_truth should be true so the design has an authoritative core")
	}
	if !card.SizesRetryAmplification {
		issues = append(issues, "sizes_retry_amplification should be true so incident behavior is reflected in capacity planning")
	}
	if !card.HasDeepDive {
		issues = append(issues, "has_deep_dive should be true so the answer demonstrates depth instead of listing components")
	}
	if !card.CoversFailureModes {
		issues = append(issues, "covers_failure_modes should be true so the design is credible under partial failure")
	}
	if !card.CoversObservability {
		issues = append(issues, "covers_observability should be true so operators can detect broken invariants")
	}
	if !card.HandlesRedesign {
		issues = append(issues, "handles_redesign should be true so changed constraints produce real design adjustments")
	}
	return issues
}

func main() {
	name := flag.String("name", "financial-consistency-drill", "scorecard name")
	flag.Parse()

	card := DrillScorecard{
		Name:                    *name,
		DefinesInvariant:        true,
		NamesSourceOfTruth:      true,
		SizesRetryAmplification: true,
		HasDeepDive:             true,
		CoversFailureModes:      true,
		CoversObservability:     true,
		HandlesRedesign:         true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateDrillScorecard(card),
	})
}
