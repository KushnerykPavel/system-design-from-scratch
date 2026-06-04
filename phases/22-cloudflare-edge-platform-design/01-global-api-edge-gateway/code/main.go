package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type GatewayConfig struct {
	PublicHost            bool   `json:"public_host"`
	TLSMode               string `json:"tls_mode"`
	RetryBudgetPercent    int    `json:"retry_budget_percent"`
	OriginShieldRegion    string `json:"origin_shield_region"`
	PolicyVersion         string `json:"policy_version"`
	AuthDependencyMode    string `json:"auth_dependency_mode"`
	ObservabilityComplete bool   `json:"observability_complete"`
}

func ValidateGatewayConfig(cfg GatewayConfig) []string {
	var issues []string
	if cfg.PublicHost && cfg.TLSMode == "" {
		issues = append(issues, "public hosts must terminate TLS explicitly")
	}
	if cfg.RetryBudgetPercent <= 0 || cfg.RetryBudgetPercent > 100 {
		issues = append(issues, "retry budget percent must be between 1 and 100")
	}
	if cfg.OriginShieldRegion == "" {
		issues = append(issues, "origin shield region should be set to reduce origin blast radius")
	}
	if cfg.PolicyVersion == "" {
		issues = append(issues, "policy version is required for rollout visibility")
	}
	if cfg.AuthDependencyMode != "local_token" && cfg.AuthDependencyMode != "remote_check" && cfg.AuthDependencyMode != "hybrid" {
		issues = append(issues, "auth dependency mode must be local_token, remote_check, or hybrid")
	}
	if !cfg.ObservabilityComplete {
		issues = append(issues, "gateway should define POP, region, and origin observability")
	}
	return issues
}

func main() {
	public := flag.Bool("public", true, "whether the listener is publicly reachable")
	flag.Parse()

	cfg := GatewayConfig{
		PublicHost:            *public,
		TLSMode:               "terminate",
		RetryBudgetPercent:    10,
		OriginShieldRegion:    "eu-central",
		PolicyVersion:         "v1",
		AuthDependencyMode:    "hybrid",
		ObservabilityComplete: true,
	}

	payload := map[string]any{
		"config": cfg,
		"issues": ValidateGatewayConfig(cfg),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
