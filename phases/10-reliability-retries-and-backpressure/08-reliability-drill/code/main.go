package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DrillAnswer struct {
	HasTimeoutOwnership  bool `json:"has_timeout_ownership"`
	HasIdempotentIngest  bool `json:"has_idempotent_ingest"`
	HasRetryBudget       bool `json:"has_retry_budget"`
	HasLoadShedding      bool `json:"has_load_shedding"`
	HasAsyncBackpressure bool `json:"has_async_backpressure"`
	HasIsolation         bool `json:"has_isolation"`
	HasObservability     bool `json:"has_observability"`
}

type DrillScore struct {
	Score   int      `json:"score"`
	Missing []string `json:"missing"`
}

func ScoreAnswer(answer DrillAnswer) DrillScore {
	score := 0
	missing := make([]string, 0, 7)

	checks := []struct {
		ok   bool
		name string
	}{
		{answer.HasTimeoutOwnership, "timeout and retry ownership"},
		{answer.HasIdempotentIngest, "idempotent ingest"},
		{answer.HasRetryBudget, "bounded retry policy"},
		{answer.HasLoadShedding, "admission control or shedding"},
		{answer.HasAsyncBackpressure, "async backpressure"},
		{answer.HasIsolation, "tenant or endpoint isolation"},
		{answer.HasObservability, "observability and operator controls"},
	}

	for _, check := range checks {
		if check.ok {
			score++
			continue
		}
		missing = append(missing, check.name)
	}

	return DrillScore{Score: score, Missing: missing}
}

func main() {
	score := ScoreAnswer(DrillAnswer{
		HasTimeoutOwnership:  true,
		HasIdempotentIngest:  true,
		HasRetryBudget:       true,
		HasLoadShedding:      true,
		HasAsyncBackpressure: true,
		HasIsolation:         true,
		HasObservability:     true,
	})

	encoded, err := json.MarshalIndent(score, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
