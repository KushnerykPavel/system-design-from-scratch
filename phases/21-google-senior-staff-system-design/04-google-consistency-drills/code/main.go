package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ConsistencyScenario struct {
	Name                     string `json:"name"`
	DefinesSourceOfTruth     bool   `json:"defines_source_of_truth"`
	NamesStrongPath          bool   `json:"names_strong_path"`
	NamesStaleReadAllowance  bool   `json:"names_stale_read_allowance"`
	DefinesAnomalyBudget     bool   `json:"defines_anomaly_budget"`
	ExplainsFailoverBehavior bool   `json:"explains_failover_behavior"`
	CoversLagMetrics         bool   `json:"covers_lag_metrics"`
	StatesTradeoffs          bool   `json:"states_tradeoffs"`
}

func ValidateConsistencyScenario(s ConsistencyScenario) []string {
	var issues []string
	if !s.DefinesSourceOfTruth {
		issues = append(issues, "defines_source_of_truth should be true so consistency claims have an authority")
	}
	if !s.NamesStrongPath {
		issues = append(issues, "names_strong_path should be true so correctness-sensitive operations are explicit")
	}
	if !s.NamesStaleReadAllowance {
		issues = append(issues, "names_stale_read_allowance should be true so weaker reads are bounded intentionally")
	}
	if !s.DefinesAnomalyBudget {
		issues = append(issues, "defines_anomaly_budget should be true so acceptable anomalies are measurable")
	}
	if !s.ExplainsFailoverBehavior {
		issues = append(issues, "explains_failover_behavior should be true so promotion and divergence risks are addressed")
	}
	if !s.CoversLagMetrics {
		issues = append(issues, "covers_lag_metrics should be true so stale data and divergence are detectable")
	}
	if !s.StatesTradeoffs {
		issues = append(issues, "states_tradeoffs should be true so stronger guarantees are not presented as free")
	}
	return issues
}

func main() {
	name := flag.String("name", "google-consistency-drill", "scenario name")
	flag.Parse()

	scenario := ConsistencyScenario{
		Name:                     *name,
		DefinesSourceOfTruth:     true,
		NamesStrongPath:          true,
		NamesStaleReadAllowance:  true,
		DefinesAnomalyBudget:     true,
		ExplainsFailoverBehavior: true,
		CoversLagMetrics:         true,
		StatesTradeoffs:          true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scenario": scenario,
		"issues":   ValidateConsistencyScenario(scenario),
	})
}
