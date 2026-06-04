package main

import (
	"encoding/json"
	"flag"
	"os"
)

type MockScorecard struct {
	Name                    string `json:"name"`
	ClarifiesPrompt         bool   `json:"clarifies_prompt"`
	PrioritizesRequirements bool   `json:"prioritizes_requirements"`
	IncludesSizing          bool   `json:"includes_sizing"`
	HasHighLevelDesign      bool   `json:"has_high_level_design"`
	HasDeepDive             bool   `json:"has_deep_dive"`
	CoversRiskAndOps        bool   `json:"covers_risk_and_ops"`
	HandlesRedesign         bool   `json:"handles_redesign"`
	StaysTimeBound          bool   `json:"stays_time_bound"`
}

func ValidateMockScorecard(card MockScorecard) []string {
	var issues []string
	if !card.ClarifiesPrompt {
		issues = append(issues, "clarifies_prompt should be true so the mock starts with a bounded problem")
	}
	if !card.PrioritizesRequirements {
		issues = append(issues, "prioritizes_requirements should be true so architecture decisions have a clear objective")
	}
	if !card.IncludesSizing {
		issues = append(issues, "includes_sizing should be true so the design is grounded in workload")
	}
	if !card.HasHighLevelDesign {
		issues = append(issues, "has_high_level_design should be true before the deep dive starts")
	}
	if !card.HasDeepDive {
		issues = append(issues, "has_deep_dive should be true so the answer demonstrates depth")
	}
	if !card.CoversRiskAndOps {
		issues = append(issues, "covers_risk_and_ops should be true so the answer includes failures, observability, and rollout")
	}
	if !card.HandlesRedesign {
		issues = append(issues, "handles_redesign should be true so changed constraints lead to a concrete update")
	}
	if !card.StaysTimeBound {
		issues = append(issues, "stays_time_bound should be true so the full loop finishes inside interview time")
	}
	return issues
}

func main() {
	name := flag.String("name", "google-full-mock", "mock name")
	flag.Parse()

	card := MockScorecard{
		Name:                    *name,
		ClarifiesPrompt:         true,
		PrioritizesRequirements: true,
		IncludesSizing:          true,
		HasHighLevelDesign:      true,
		HasDeepDive:             true,
		CoversRiskAndOps:        true,
		HandlesRedesign:         true,
		StaysTimeBound:          true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateMockScorecard(card),
	})
}
