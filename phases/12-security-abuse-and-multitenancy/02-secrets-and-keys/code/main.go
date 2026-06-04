package main

import (
	"encoding/json"
	"os"
)

type RotationPlan struct {
	ShortLivedCredentials bool `json:"short_lived_credentials"`
	BackgroundRefresh     bool `json:"background_refresh"`
	OverlapWindow         bool `json:"overlap_window"`
	RevocationPath        bool `json:"revocation_path"`
	HotPathSecretFetch    bool `json:"hot_path_secret_fetch"`
}

func ValidateRotationPlan(plan RotationPlan) []string {
	var issues []string
	if !plan.ShortLivedCredentials {
		issues = append(issues, "prefer short-lived credentials where possible")
	}
	if !plan.BackgroundRefresh {
		issues = append(issues, "rotation should refresh in background rather than require restarts")
	}
	if !plan.OverlapWindow {
		issues = append(issues, "rotation should support overlap windows or dual trust")
	}
	if !plan.RevocationPath {
		issues = append(issues, "compromised credentials need an explicit revocation path")
	}
	if plan.HotPathSecretFetch {
		issues = append(issues, "request path should not synchronously depend on secret fetches")
	}
	return issues
}

func main() {
	plan := RotationPlan{
		ShortLivedCredentials: true,
		BackgroundRefresh:     true,
		OverlapWindow:         true,
		RevocationPath:        true,
		HotPathSecretFetch:    false,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateRotationPlan(plan),
	})
}
