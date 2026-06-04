package main

import (
	"encoding/json"
	"flag"
	"os"
)

type AlertPolicy struct {
	Name                 string `json:"name"`
	EscalationLevels     int    `json:"escalation_levels"`
	DedupWindowMinutes   int    `json:"dedup_window_minutes"`
	AckTimeoutMinutes    int    `json:"ack_timeout_minutes"`
	OwnershipRequired    bool   `json:"ownership_required"`
	RunbookRequired      bool   `json:"runbook_required"`
	ProviderFallback     bool   `json:"provider_fallback"`
	MaintenanceSupported bool   `json:"maintenance_supported"`
}

func ValidateAlertPolicy(policy AlertPolicy) []string {
	var issues []string
	if policy.EscalationLevels < 1 {
		issues = append(issues, "escalation_levels should be at least 1")
	}
	if policy.DedupWindowMinutes < 1 {
		issues = append(issues, "dedup_window_minutes should be positive")
	}
	if policy.AckTimeoutMinutes < 1 {
		issues = append(issues, "ack_timeout_minutes should be positive")
	}
	if !policy.OwnershipRequired {
		issues = append(issues, "ownership_required should be true for routable sev-1 alerts")
	}
	if !policy.RunbookRequired {
		issues = append(issues, "runbook_required should usually be true for production paging")
	}
	if !policy.ProviderFallback {
		issues = append(issues, "provider_fallback should reduce single-vendor paging risk")
	}
	if !policy.MaintenanceSupported {
		issues = append(issues, "maintenance_supported should exist to avoid noisy known work")
	}
	return issues
}

func main() {
	name := flag.String("name", "default-production", "alert policy name")
	flag.Parse()

	policy := AlertPolicy{
		Name:                 *name,
		EscalationLevels:     2,
		DedupWindowMinutes:   5,
		AckTimeoutMinutes:    10,
		OwnershipRequired:    true,
		RunbookRequired:      true,
		ProviderFallback:     true,
		MaintenanceSupported: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"policy": policy,
		"issues": ValidateAlertPolicy(policy),
	})
}
