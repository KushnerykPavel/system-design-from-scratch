package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type LimiterConfig struct {
	Name               string `json:"name"`
	Mode               string `json:"mode"`
	Store              string `json:"store"`
	LocalTTLMillis     int    `json:"local_ttl_millis"`
	Burst              int    `json:"burst"`
	RetryBudgetPercent int    `json:"retry_budget_percent"`
	FailOpen           bool   `json:"fail_open"`
	SharedEnforcement  bool   `json:"shared_enforcement"`
	Explainability     bool   `json:"explainability"`
}

func ValidateConfig(cfg LimiterConfig) []string {
	var issues []string
	if cfg.Mode != "token_bucket" && cfg.Mode != "sliding_window" {
		issues = append(issues, "mode must be token_bucket or sliding_window")
	}
	if cfg.Store == "" {
		issues = append(issues, "shared store is required")
	}
	if cfg.Burst <= 0 {
		issues = append(issues, "burst must be positive")
	}
	if cfg.SharedEnforcement && cfg.LocalTTLMillis <= 0 {
		issues = append(issues, "shared enforcement requires a positive local_ttl_millis")
	}
	if cfg.FailOpen && cfg.RetryBudgetPercent == 0 {
		issues = append(issues, "fail-open mode should define a retry budget or bounded degraded mode")
	}
	if !cfg.Explainability {
		issues = append(issues, "reject decisions should be explainable for operators and customers")
	}
	return issues
}

func main() {
	name := flag.String("name", "edge-api-limiter", "name of the limiter config")
	flag.Parse()

	cfg := LimiterConfig{
		Name:               *name,
		Mode:               "token_bucket",
		Store:              "redis",
		LocalTTLMillis:     100,
		Burst:              50,
		RetryBudgetPercent: 10,
		FailOpen:           true,
		SharedEnforcement:  true,
		Explainability:     true,
	}

	payload := map[string]any{
		"config": cfg,
		"issues": ValidateConfig(cfg),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
