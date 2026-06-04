package main

import (
	"encoding/json"
	"flag"
	"os"
)

type SchedulerPlan struct {
	Name                 string `json:"name"`
	HasShardLeases       bool   `json:"has_shard_leases"`
	HasRetryPolicy       bool   `json:"has_retry_policy"`
	HasJitter            bool   `json:"has_jitter"`
	HandlesMissedRuns    bool   `json:"handles_missed_runs"`
	HasTenantQuotas      bool   `json:"has_tenant_quotas"`
	HasDLQ               bool   `json:"has_dlq"`
	HasDeadlineAwareness bool   `json:"has_deadline_awareness"`
}

func ValidateSchedulerPlan(plan SchedulerPlan) []string {
	var issues []string
	if !plan.HasShardLeases {
		issues = append(issues, "has_shard_leases should be true so scheduler ownership recovers after failure")
	}
	if !plan.HasRetryPolicy {
		issues = append(issues, "has_retry_policy should be true so failed jobs have bounded behavior")
	}
	if !plan.HasJitter {
		issues = append(issues, "has_jitter should be true to avoid synchronized retry or cron storms")
	}
	if !plan.HandlesMissedRuns {
		issues = append(issues, "handles_missed_runs should be true so downtime does not create ambiguous catch-up behavior")
	}
	if !plan.HasTenantQuotas {
		issues = append(issues, "has_tenant_quotas should be true so one tenant cannot consume all scheduler capacity")
	}
	if !plan.HasDLQ {
		issues = append(issues, "has_dlq should be true for jobs that repeatedly fail or cannot be retried safely")
	}
	if !plan.HasDeadlineAwareness {
		issues = append(issues, "has_deadline_awareness should be true so stale jobs can be dropped or deprioritized")
	}
	return issues
}

func main() {
	name := flag.String("name", "job-scheduler", "plan name")
	flag.Parse()

	plan := SchedulerPlan{
		Name:                 *name,
		HasShardLeases:       true,
		HasRetryPolicy:       true,
		HasJitter:            true,
		HandlesMissedRuns:    true,
		HasTenantQuotas:      true,
		HasDLQ:               true,
		HasDeadlineAwareness: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateSchedulerPlan(plan),
	})
}
