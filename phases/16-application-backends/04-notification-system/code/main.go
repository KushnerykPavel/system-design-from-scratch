package main

import (
	"encoding/json"
	"flag"
	"os"
)

type NotificationPlan struct {
	Name                     string `json:"name"`
	HasDedupeStore           bool   `json:"has_dedupe_store"`
	HasPreferenceChecks      bool   `json:"has_preference_checks"`
	UsesPriorityQueues       bool   `json:"uses_priority_queues"`
	HasRetryBudget           bool   `json:"has_retry_budget"`
	SupportsProviderFailover bool   `json:"supports_provider_failover"`
	QuietHoursEnforced       bool   `json:"quiet_hours_enforced"`
	TracksDeliveryState      bool   `json:"tracks_delivery_state"`
}

func ValidateNotificationPlan(plan NotificationPlan) []string {
	var issues []string
	if !plan.HasDedupeStore {
		issues = append(issues, "has_dedupe_store should be true to avoid duplicate sends")
	}
	if !plan.HasPreferenceChecks {
		issues = append(issues, "has_preference_checks should be true to enforce opt-outs and policy rules")
	}
	if !plan.UsesPriorityQueues {
		issues = append(issues, "uses_priority_queues should be true so low-priority traffic does not crowd out critical messages")
	}
	if !plan.HasRetryBudget {
		issues = append(issues, "has_retry_budget should be true to prevent infinite retry storms")
	}
	if !plan.SupportsProviderFailover {
		issues = append(issues, "supports_provider_failover should be true for critical notification classes")
	}
	if !plan.QuietHoursEnforced {
		issues = append(issues, "quiet_hours_enforced should be true when the product exposes quiet-hour controls")
	}
	if !plan.TracksDeliveryState {
		issues = append(issues, "tracks_delivery_state should be true for auditability and retries")
	}
	return issues
}

func main() {
	name := flag.String("name", "notification-platform", "plan name")
	flag.Parse()

	plan := NotificationPlan{
		Name:                     *name,
		HasDedupeStore:           true,
		HasPreferenceChecks:      true,
		UsesPriorityQueues:       true,
		HasRetryBudget:           true,
		SupportsProviderFailover: true,
		QuietHoursEnforced:       true,
		TracksDeliveryState:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateNotificationPlan(plan),
	})
}
