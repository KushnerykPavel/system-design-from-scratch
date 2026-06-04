package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type MockScorecard struct {
	ScopeClarity            bool `json:"scope_clarity"`
	SizingIncluded          bool `json:"sizing_included"`
	ControlDataPlaneSplit   bool `json:"control_data_plane_split"`
	DeepDiveIncluded        bool `json:"deep_dive_included"`
	FailureCoverage         bool `json:"failure_coverage"`
	ObservabilityCoverage   bool `json:"observability_coverage"`
	RedesignIncluded        bool `json:"redesign_included"`
	AmplificationRiskNamed  bool `json:"amplification_risk_named"`
	ExplainabilityAddressed bool `json:"explainability_addressed"`
}

func ValidateMockScorecard(card MockScorecard) []string {
	var issues []string
	if !card.ScopeClarity {
		issues = append(issues, "mock answer should clarify scope and edge priorities")
	}
	if !card.SizingIncluded {
		issues = append(issues, "mock answer should include at least one useful estimate")
	}
	if !card.ControlDataPlaneSplit {
		issues = append(issues, "answer should separate control plane and data plane")
	}
	if !card.DeepDiveIncluded {
		issues = append(issues, "answer should include one deliberate deep dive")
	}
	if !card.FailureCoverage || !card.ObservabilityCoverage {
		issues = append(issues, "answer should cover failure modes and observability")
	}
	if !card.RedesignIncluded {
		issues = append(issues, "answer should include redesign under changed assumptions")
	}
	if !card.AmplificationRiskNamed {
		issues = append(issues, "answer should name at least one amplification risk")
	}
	if !card.ExplainabilityAddressed {
		issues = append(issues, "answer should preserve operator or customer explainability")
	}
	return issues
}

func main() {
	mode := flag.String("mode", "cloudflare", "mock mode")
	flag.Parse()

	card := MockScorecard{
		ScopeClarity:            true,
		SizingIncluded:          true,
		ControlDataPlaneSplit:   true,
		DeepDiveIncluded:        true,
		FailureCoverage:         true,
		ObservabilityCoverage:   true,
		RedesignIncluded:        true,
		AmplificationRiskNamed:  true,
		ExplainabilityAddressed: true,
	}

	payload := map[string]any{
		"mode":   *mode,
		"card":   card,
		"issues": ValidateMockScorecard(card),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
