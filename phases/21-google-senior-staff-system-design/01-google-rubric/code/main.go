package main

import (
	"encoding/json"
	"flag"
	"os"
)

type RubricScorecard struct {
	Name                   string `json:"name"`
	ClarifiesScope         bool   `json:"clarifies_scope"`
	QuantifiesWorkload     bool   `json:"quantifies_workload"`
	HasHighLevelDesign     bool   `json:"has_high_level_design"`
	ChoosesDeepDive        bool   `json:"chooses_deep_dive"`
	ExplainsTradeoffs      bool   `json:"explains_tradeoffs"`
	CoversFailureModes     bool   `json:"covers_failure_modes"`
	CoversObservability    bool   `json:"covers_observability"`
	HandlesRedesignCleanly bool   `json:"handles_redesign_cleanly"`
}

func ValidateRubricScorecard(card RubricScorecard) []string {
	var issues []string
	if !card.ClarifiesScope {
		issues = append(issues, "clarifies_scope should be true so the answer starts with the real problem")
	}
	if !card.QuantifiesWorkload {
		issues = append(issues, "quantifies_workload should be true so the architecture is grounded in scale")
	}
	if !card.HasHighLevelDesign {
		issues = append(issues, "has_high_level_design should be true before any deep dive")
	}
	if !card.ChoosesDeepDive {
		issues = append(issues, "chooses_deep_dive should be true so depth is deliberate")
	}
	if !card.ExplainsTradeoffs {
		issues = append(issues, "explains_tradeoffs should be true so the answer shows judgment")
	}
	if !card.CoversFailureModes {
		issues = append(issues, "covers_failure_modes should be true so the design is credible under stress")
	}
	if !card.CoversObservability {
		issues = append(issues, "covers_observability should be true so operators can see problems")
	}
	if !card.HandlesRedesignCleanly {
		issues = append(issues, "handles_redesign_cleanly should be true so changed constraints lead to concrete updates")
	}
	return issues
}

func main() {
	name := flag.String("name", "google-rubric", "scorecard name")
	flag.Parse()

	card := RubricScorecard{
		Name:                   *name,
		ClarifiesScope:         true,
		QuantifiesWorkload:     true,
		HasHighLevelDesign:     true,
		ChoosesDeepDive:        true,
		ExplainsTradeoffs:      true,
		CoversFailureModes:     true,
		CoversObservability:    true,
		HandlesRedesignCleanly: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateRubricScorecard(card),
	})
}
