package main

import (
	"encoding/json"
	"flag"
	"os"
)

type PanelReview struct {
	Name          string `json:"name"`
	Clarification int    `json:"clarification"`
	Sizing        int    `json:"sizing"`
	Architecture  int    `json:"architecture"`
	DeepDive      int    `json:"deep_dive"`
	FailureModes  int    `json:"failure_modes"`
	Observability int    `json:"observability"`
	TradeOffs     int    `json:"trade_offs"`
	Communication int    `json:"communication"`
}

func ValidatePanelReview(review PanelReview) []string {
	var issues []string
	for field, score := range map[string]int{
		"clarification": review.Clarification,
		"sizing":        review.Sizing,
		"architecture":  review.Architecture,
		"deep_dive":     review.DeepDive,
		"failure_modes": review.FailureModes,
		"observability": review.Observability,
		"trade_offs":    review.TradeOffs,
		"communication": review.Communication,
	} {
		if score < 1 || score > 4 {
			issues = append(issues, field+" score must be between 1 and 4")
		}
	}
	return issues
}

func main() {
	name := flag.String("name", "final-capstone-review", "review name")
	flag.Parse()

	review := PanelReview{
		Name:          *name,
		Clarification: 4,
		Sizing:        4,
		Architecture:  4,
		DeepDive:      4,
		FailureModes:  4,
		Observability: 4,
		TradeOffs:     4,
		Communication: 4,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"review": review,
		"issues": ValidatePanelReview(review),
	})
}
