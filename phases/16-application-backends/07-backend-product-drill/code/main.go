package main

import (
	"encoding/json"
	"flag"
	"os"
)

type AnswerScorecard struct {
	Clarifications   int  `json:"clarifications"`
	HasSizing        bool `json:"has_sizing"`
	HasRequirements  bool `json:"has_requirements"`
	HasArchitecture  bool `json:"has_architecture"`
	DeepDives        int  `json:"deep_dives"`
	HasFailureModes  bool `json:"has_failure_modes"`
	HasObservability bool `json:"has_observability"`
	HasTradeoffs     bool `json:"has_tradeoffs"`
	HasRollout       bool `json:"has_rollout"`
}

func ScoreAnswer(card AnswerScorecard) []string {
	var gaps []string
	if card.Clarifications < 2 {
		gaps = append(gaps, "ask at least two clarifying questions before committing to the design")
	}
	if !card.HasSizing {
		gaps = append(gaps, "add rough sizing before architecture choices")
	}
	if !card.HasRequirements {
		gaps = append(gaps, "state functional and non-functional requirements explicitly")
	}
	if !card.HasArchitecture {
		gaps = append(gaps, "include a high-level architecture")
	}
	if card.DeepDives < 1 {
		gaps = append(gaps, "choose at least one intentional deep dive")
	}
	if !card.HasFailureModes {
		gaps = append(gaps, "add failure modes and mitigations")
	}
	if !card.HasObservability {
		gaps = append(gaps, "attach metrics or SLOs to the main product promise")
	}
	if !card.HasTradeoffs {
		gaps = append(gaps, "name explicit trade-offs instead of only listing components")
	}
	if !card.HasRollout {
		gaps = append(gaps, "include rollout or migration thinking")
	}
	return gaps
}

func main() {
	name := flag.String("name", "phase-16-drill", "scorecard name")
	flag.Parse()

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"name": *name,
		"gaps": ScoreAnswer(AnswerScorecard{
			Clarifications:   3,
			HasSizing:        true,
			HasRequirements:  true,
			HasArchitecture:  true,
			DeepDives:        2,
			HasFailureModes:  true,
			HasObservability: true,
			HasTradeoffs:     true,
			HasRollout:       true,
		}),
	})
}
