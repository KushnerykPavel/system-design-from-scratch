package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type FailoverConfig struct {
	PoolName             string `json:"pool_name"`
	HealthMode           string `json:"health_mode"`
	RetryBudgetPercent   int    `json:"retry_budget_percent"`
	CooldownSeconds      int    `json:"cooldown_seconds"`
	ShieldEnabled        bool   `json:"shield_enabled"`
	CrossRegionFailover  bool   `json:"cross_region_failover"`
	Explainability       bool   `json:"explainability"`
	ReserveCapacityRatio int    `json:"reserve_capacity_ratio"`
}

func ValidateFailoverConfig(cfg FailoverConfig) []string {
	var issues []string
	if cfg.PoolName == "" {
		issues = append(issues, "pool name is required")
	}
	if cfg.HealthMode != "binary" && cfg.HealthMode != "scored" && cfg.HealthMode != "hybrid" {
		issues = append(issues, "health mode must be binary, scored, or hybrid")
	}
	if cfg.RetryBudgetPercent <= 0 || cfg.RetryBudgetPercent > 100 {
		issues = append(issues, "retry budget percent must be between 1 and 100")
	}
	if cfg.CooldownSeconds <= 0 {
		issues = append(issues, "cooldown seconds must be positive")
	}
	if cfg.CrossRegionFailover && cfg.ReserveCapacityRatio < 120 {
		issues = append(issues, "cross-region failover should reserve at least 120 percent capacity")
	}
	if !cfg.ShieldEnabled {
		issues = append(issues, "origin shielding should be considered for fragile or shared origins")
	}
	if !cfg.Explainability {
		issues = append(issues, "origin selection and failover reasons should be explainable")
	}
	return issues
}

func main() {
	pool := flag.String("pool", "primary-api-pool", "origin pool name")
	flag.Parse()

	cfg := FailoverConfig{
		PoolName:             *pool,
		HealthMode:           "hybrid",
		RetryBudgetPercent:   10,
		CooldownSeconds:      30,
		ShieldEnabled:        true,
		CrossRegionFailover:  true,
		Explainability:       true,
		ReserveCapacityRatio: 150,
	}

	payload := map[string]any{
		"config": cfg,
		"issues": ValidateFailoverConfig(cfg),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
