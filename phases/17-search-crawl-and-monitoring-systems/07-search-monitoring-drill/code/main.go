package main

import (
	"encoding/json"
	"flag"
	"os"
)

type DrillScorecard struct {
	Name            string `json:"name"`
	ClarifiedScope  bool   `json:"clarified_scope"`
	SizedWorkload   bool   `json:"sized_workload"`
	ChosenDeepDive  bool   `json:"chosen_deep_dive"`
	FailureModes    bool   `json:"failure_modes"`
	Observability   bool   `json:"observability"`
	Tradeoffs       bool   `json:"tradeoffs"`
	RedesignHandled bool   `json:"redesign_handled"`
}

func ValidateDrillScorecard(card DrillScorecard) []string {
	var issues []string
	if !card.ClarifiedScope {
		issues = append(issues, "clarified_scope should be true before architecture")
	}
	if !card.SizedWorkload {
		issues = append(issues, "sized_workload should be true before detailed design claims")
	}
	if !card.ChosenDeepDive {
		issues = append(issues, "chosen_deep_dive should be true for a senior-level drill")
	}
	if !card.FailureModes {
		issues = append(issues, "failure_modes should be covered explicitly")
	}
	if !card.Observability {
		issues = append(issues, "observability should include concrete signals")
	}
	if !card.Tradeoffs {
		issues = append(issues, "tradeoffs should be named with benefit and cost")
	}
	if !card.RedesignHandled {
		issues = append(issues, "redesign_handled should be true after a constraint change")
	}
	return issues
}

func main() {
	name := flag.String("name", "phase-17-drill", "drill scorecard name")
	flag.Parse()

	card := DrillScorecard{
		Name:            *name,
		ClarifiedScope:  true,
		SizedWorkload:   true,
		ChosenDeepDive:  true,
		FailureModes:    true,
		Observability:   true,
		Tradeoffs:       true,
		RedesignHandled: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateDrillScorecard(card),
	})
}
