package main

import (
	"encoding/json"
	"flag"
	"os"
)

type IndexUpdatePlan struct {
	Name                    string `json:"name"`
	DocumentLagSeconds      int    `json:"document_lag_seconds"`
	RankingLagSeconds       int    `json:"ranking_lag_seconds"`
	BlueGreenEnabled        bool   `json:"blue_green_enabled"`
	DualReadValidation      bool   `json:"dual_read_validation"`
	DeleteTombstonesEnabled bool   `json:"delete_tombstones_enabled"`
	BackfillWorkerPools     int    `json:"backfill_worker_pools"`
	SchemaCompatChecks      bool   `json:"schema_compat_checks"`
}

func ValidateIndexUpdatePlan(plan IndexUpdatePlan) []string {
	var issues []string
	if plan.DocumentLagSeconds <= 0 {
		issues = append(issues, "document_lag_seconds must be positive")
	}
	if plan.RankingLagSeconds <= 0 {
		issues = append(issues, "ranking_lag_seconds must be positive")
	}
	if !plan.BlueGreenEnabled {
		issues = append(issues, "blue_green_enabled should be true for safe publish rollback")
	}
	if !plan.DualReadValidation {
		issues = append(issues, "dual_read_validation should be enabled for canary result checks")
	}
	if !plan.DeleteTombstonesEnabled {
		issues = append(issues, "delete_tombstones_enabled should be true for fast removals")
	}
	if plan.BackfillWorkerPools < 1 {
		issues = append(issues, "backfill_worker_pools should isolate live freshness from backfills")
	}
	if !plan.SchemaCompatChecks {
		issues = append(issues, "schema_compat_checks should be enabled before rollout")
	}
	return issues
}

func main() {
	name := flag.String("name", "search-index-updates", "index update plan name")
	flag.Parse()

	plan := IndexUpdatePlan{
		Name:                    *name,
		DocumentLagSeconds:      300,
		RankingLagSeconds:       900,
		BlueGreenEnabled:        true,
		DualReadValidation:      true,
		DeleteTombstonesEnabled: true,
		BackfillWorkerPools:     2,
		SchemaCompatChecks:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateIndexUpdatePlan(plan),
	})
}
