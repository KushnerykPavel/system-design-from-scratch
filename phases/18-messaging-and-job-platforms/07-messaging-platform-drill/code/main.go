package main

import (
	"encoding/json"
	"flag"
	"os"
)

type DrillScorecard struct {
	Name                 string `json:"name"`
	ClarifiedScope       bool   `json:"clarified_scope"`
	SizedWorkload        bool   `json:"sized_workload"`
	ExplainedGuarantees  bool   `json:"explained_guarantees"`
	ChosenDeepDive       bool   `json:"chosen_deep_dive"`
	CoveredFailures      bool   `json:"covered_failures"`
	CoveredObservability bool   `json:"covered_observability"`
	HandledRedesign      bool   `json:"handled_redesign"`
}

func ValidateDrillScorecard(card DrillScorecard) []string {
	var issues []string
	if !card.ClarifiedScope {
		issues = append(issues, "clarified_scope should be true before architecture")
	}
	if !card.SizedWorkload {
		issues = append(issues, "sized_workload should be true before scale claims")
	}
	if !card.ExplainedGuarantees {
		issues = append(issues, "explained_guarantees should be true for messaging semantics")
	}
	if !card.ChosenDeepDive {
		issues = append(issues, "chosen_deep_dive should be true for a senior-level drill")
	}
	if !card.CoveredFailures {
		issues = append(issues, "covered_failures should be true with concrete detection and mitigation")
	}
	if !card.CoveredObservability {
		issues = append(issues, "covered_observability should be true with queue or workflow signals")
	}
	if !card.HandledRedesign {
		issues = append(issues, "handled_redesign should be true after the prompt changes")
	}
	return issues
}

func main() {
	name := flag.String("name", "messaging-platform-drill", "drill name")
	flag.Parse()

	card := DrillScorecard{
		Name:                 *name,
		ClarifiedScope:       true,
		SizedWorkload:        true,
		ExplainedGuarantees:  true,
		ChosenDeepDive:       true,
		CoveredFailures:      true,
		CoveredObservability: true,
		HandledRedesign:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateDrillScorecard(card),
	})
}
