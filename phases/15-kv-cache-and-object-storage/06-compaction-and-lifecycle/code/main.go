package main

import (
	"encoding/json"
	"flag"
	"os"
)

type LifecyclePolicy struct {
	Name                string `json:"name"`
	TransitionAfterDays int    `json:"transition_after_days"`
	ExpireAfterDays     int    `json:"expire_after_days"`
	DeleteGraceDays     int    `json:"delete_grace_days"`
	LegalHoldAware      bool   `json:"legal_hold_aware"`
	DryRunRequired      bool   `json:"dry_run_required"`
	CompactionBudgetPct int    `json:"compaction_budget_pct"`
	GCGraceHours        int    `json:"gc_grace_hours"`
}

func ValidateLifecyclePolicy(cfg LifecyclePolicy) []string {
	var issues []string
	if cfg.TransitionAfterDays < 0 {
		issues = append(issues, "transition_after_days cannot be negative")
	}
	if cfg.ExpireAfterDays > 0 && cfg.TransitionAfterDays > 0 && cfg.ExpireAfterDays <= cfg.TransitionAfterDays {
		issues = append(issues, "expire_after_days should be greater than transition_after_days")
	}
	if cfg.DeleteGraceDays <= 0 {
		issues = append(issues, "delete_grace_days must be positive")
	}
	if !cfg.LegalHoldAware {
		issues = append(issues, "legal_hold_aware should be true for destructive lifecycle actions")
	}
	if !cfg.DryRunRequired {
		issues = append(issues, "dry_run_required should be true before activating destructive rules")
	}
	if cfg.CompactionBudgetPct < 10 || cfg.CompactionBudgetPct > 50 {
		issues = append(issues, "compaction_budget_pct should stay between 10 and 50")
	}
	if cfg.GCGraceHours <= 0 {
		issues = append(issues, "gc_grace_hours must be positive")
	}
	return issues
}

func main() {
	name := flag.String("name", "default-lifecycle", "name of the lifecycle policy")
	flag.Parse()

	cfg := LifecyclePolicy{
		Name:                *name,
		TransitionAfterDays: 30,
		ExpireAfterDays:     365,
		DeleteGraceDays:     7,
		LegalHoldAware:      true,
		DryRunRequired:      true,
		CompactionBudgetPct: 25,
		GCGraceHours:        72,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"policy": cfg,
		"issues": ValidateLifecyclePolicy(cfg),
	})
}
