package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BudgetInput struct {
	TotalEvents int64   `json:"total_events"`
	BadEvents   int64   `json:"bad_events"`
	TargetRatio float64 `json:"target_ratio"`
	WindowDays  int     `json:"window_days"`
	ElapsedDays int     `json:"elapsed_days"`
}

type BudgetAssessment struct {
	AllowedBadEvents int64   `json:"allowed_bad_events"`
	ConsumedRatio    float64 `json:"consumed_ratio"`
	BurnRate         float64 `json:"burn_rate"`
	Status           string  `json:"status"`
}

func AssessBudget(in BudgetInput) BudgetAssessment {
	if in.TargetRatio <= 0 || in.TargetRatio >= 1 || in.TotalEvents <= 0 || in.WindowDays <= 0 || in.ElapsedDays <= 0 {
		return BudgetAssessment{Status: "invalid"}
	}

	allowed := int64((1 - in.TargetRatio) * float64(in.TotalEvents))
	if allowed < 1 {
		allowed = 1
	}

	consumed := float64(in.BadEvents) / float64(allowed)
	elapsed := float64(in.ElapsedDays) / float64(in.WindowDays)
	burnRate := consumed / elapsed

	status := "healthy"
	switch {
	case consumed >= 1 || burnRate >= 4:
		status = "critical"
	case burnRate >= 2:
		status = "warning"
	}

	return BudgetAssessment{
		AllowedBadEvents: allowed,
		ConsumedRatio:    consumed,
		BurnRate:         burnRate,
		Status:           status,
	}
}

func main() {
	input := BudgetInput{
		TotalEvents: 2000000000,
		BadEvents:   700000,
		TargetRatio: 0.999,
		WindowDays:  30,
		ElapsedDays: 10,
	}

	encoded, err := json.MarshalIndent(AssessBudget(input), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
