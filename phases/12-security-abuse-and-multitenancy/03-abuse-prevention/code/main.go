package main

import (
	"encoding/json"
	"os"
)

type AbusePlan struct {
	HasEdgeLimits       bool `json:"has_edge_limits"`
	HasAccountAwareGate bool `json:"has_account_aware_gate"`
	HasTenantBudget     bool `json:"has_tenant_budget"`
	HasChallengeStep    bool `json:"has_challenge_step"`
	IPOnlyEnforcement   bool `json:"ip_only_enforcement"`
}

func ValidateAbusePlan(plan AbusePlan) []string {
	var issues []string
	if !plan.HasEdgeLimits {
		issues = append(issues, "start with cheap edge filtering or coarse rate limits")
	}
	if !plan.HasAccountAwareGate {
		issues = append(issues, "sensitive flows need account or API-key aware controls")
	}
	if !plan.HasTenantBudget {
		issues = append(issues, "shared infrastructure should include tenant-aware protection")
	}
	if !plan.HasChallengeStep {
		issues = append(issues, "consider challenge or friction before hard blocking healthy-looking traffic")
	}
	if plan.IPOnlyEnforcement {
		issues = append(issues, "IP-only enforcement is weak under NATs and adversarial churn")
	}
	return issues
}

func main() {
	plan := AbusePlan{
		HasEdgeLimits:       true,
		HasAccountAwareGate: true,
		HasTenantBudget:     true,
		HasChallengeStep:    true,
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateAbusePlan(plan),
	})
}
