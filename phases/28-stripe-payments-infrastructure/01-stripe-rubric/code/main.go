package main

import (
	"encoding/json"
	"os"
)

type DesignSignal struct {
	Category    string
	Description string
	Observed    bool
}

type RubricResult struct {
	TotalSignals    int      `json:"total_signals"`
	ObservedSignals int      `json:"observed_signals"`
	ScorePercent    int      `json:"score_percent"`
	HireSignal      string   `json:"hire_signal"`
	MissingAreas    []string `json:"missing_areas"`
}

func DefaultSignals() []DesignSignal {
	return []DesignSignal{
		{Category: "idempotency", Description: "designed idempotency keys on every mutating operation", Observed: false},
		{Category: "double-entry", Description: "used double-entry ledger with debit/credit pairs", Observed: false},
		{Category: "atomicity", Description: "charge either completes fully or reverses fully — no partial state", Observed: false},
		{Category: "pci-compliance", Description: "isolated cardholder data with tokenization", Observed: false},
		{Category: "fraud-hooks", Description: "integrated fraud scoring without blocking critical path", Observed: false},
		{Category: "reconciliation", Description: "designed reconciliation against external settlement files", Observed: false},
		{Category: "webhook-delivery", Description: "designed reliable webhook delivery with retry and DLQ", Observed: false},
		{Category: "tradeoffs", Description: "named trade-offs for consistency vs availability in payment path", Observed: false},
	}
}

func EvaluateRubric(signals []DesignSignal) RubricResult {
	total := len(signals)
	observed := 0
	var missing []string
	for _, s := range signals {
		if s.Observed {
			observed++
		} else {
			missing = append(missing, s.Category+": "+s.Description)
		}
	}
	score := 0
	if total > 0 {
		score = (observed * 100) / total
	}
	hire := "no-hire"
	switch {
	case score >= 88:
		hire = "strong-hire"
	case score >= 63:
		hire = "hire"
	case score >= 38:
		hire = "mixed"
	}
	return RubricResult{
		TotalSignals:    total,
		ObservedSignals: observed,
		ScorePercent:    score,
		HireSignal:      hire,
		MissingAreas:    missing,
	}
}

func main() {
	signals := DefaultSignals()
	for i := range signals {
		if signals[i].Category != "pci-compliance" {
			signals[i].Observed = true
		}
	}
	result := EvaluateRubric(signals)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
