package main

import (
	"encoding/json"
	"os"
)

type IsolationPlan struct {
	TenantScopedKeys  bool `json:"tenant_scoped_keys"`
	FairScheduling    bool `json:"fair_scheduling"`
	PerTenantBudgets  bool `json:"per_tenant_budgets"`
	CarveOutPath      bool `json:"carve_out_path"`
	SharedAdminBypass bool `json:"shared_admin_bypass"`
}

func ValidateIsolationPlan(plan IsolationPlan) []string {
	var issues []string
	if !plan.TenantScopedKeys {
		issues = append(issues, "shared state should be tenant-scoped explicitly")
	}
	if !plan.FairScheduling {
		issues = append(issues, "shared workers need fair scheduling or bulkheads")
	}
	if !plan.PerTenantBudgets {
		issues = append(issues, "performance isolation usually requires per-tenant budgets")
	}
	if !plan.CarveOutPath {
		issues = append(issues, "heavy tenants need a migration or carve-out path")
	}
	if plan.SharedAdminBypass {
		issues = append(issues, "broad admin bypass weakens tenant isolation")
	}
	return issues
}

func main() {
	plan := IsolationPlan{
		TenantScopedKeys: true,
		FairScheduling:   true,
		PerTenantBudgets: true,
		CarveOutPath:     true,
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateIsolationPlan(plan),
	})
}
