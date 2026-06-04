package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BudgetPolicy struct {
	BaseQPS              int     `json:"base_qps"`
	RetryBudgetRatio     float64 `json:"retry_budget_ratio"`
	HedgeBudgetRatio     float64 `json:"hedge_budget_ratio"`
	HedgeAfterMS         int     `json:"hedge_after_ms"`
	SupportsCancellation bool    `json:"supports_cancellation"`
	SafeToHedge          bool    `json:"safe_to_hedge"`
}

type BudgetAssessment struct {
	AllowedExtraAttempts int      `json:"allowed_extra_attempts"`
	Risk                 string   `json:"risk"`
	Notes                []string `json:"notes"`
}

func AssessBudgetPolicy(policy BudgetPolicy) BudgetAssessment {
	notes := make([]string, 0, 4)
	score := 0
	totalBudget := policy.RetryBudgetRatio + policy.HedgeBudgetRatio
	allowed := int(float64(policy.BaseQPS) * totalBudget)

	if totalBudget > 0.12 {
		score += 2
		notes = append(notes, "speculative budget is large relative to baseline traffic")
	}
	if policy.HedgeBudgetRatio > 0 && !policy.SafeToHedge {
		score += 2
		notes = append(notes, "hedging is enabled for unsafe or stateful requests")
	}
	if policy.HedgeBudgetRatio > 0 && !policy.SupportsCancellation {
		score++
		notes = append(notes, "losing hedged attempts cannot be canceled")
	}
	if policy.HedgeBudgetRatio > 0 && policy.HedgeAfterMS == 0 {
		score++
		notes = append(notes, "hedges fire immediately instead of targeting tail latency")
	}

	risk := "low"
	if score >= 4 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return BudgetAssessment{
		AllowedExtraAttempts: allowed,
		Risk:                 risk,
		Notes:                notes,
	}
}

func main() {
	assessment := AssessBudgetPolicy(BudgetPolicy{
		BaseQPS:              180000,
		RetryBudgetRatio:     0.05,
		HedgeBudgetRatio:     0.02,
		HedgeAfterMS:         120,
		SupportsCancellation: true,
		SafeToHedge:          true,
	})

	encoded, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
