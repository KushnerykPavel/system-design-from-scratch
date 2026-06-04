package main

import (
	"encoding/json"
	"flag"
	"os"
)

type PubSubPlan struct {
	Name                    string `json:"name"`
	HasTopicPartitions      bool   `json:"has_topic_partitions"`
	HasPerSubscriptionState bool   `json:"has_per_subscription_state"`
	SupportsFiltering       bool   `json:"supports_filtering"`
	HasReplayCursor         bool   `json:"has_replay_cursor"`
	HasIsolationControls    bool   `json:"has_isolation_controls"`
	HasBacklogBudgets       bool   `json:"has_backlog_budgets"`
	HasDeadLetterPath       bool   `json:"has_dead_letter_path"`
}

func ValidatePubSubPlan(plan PubSubPlan) []string {
	var issues []string
	if !plan.HasTopicPartitions {
		issues = append(issues, "has_topic_partitions should be true so fanout throughput can scale")
	}
	if !plan.HasPerSubscriptionState {
		issues = append(issues, "has_per_subscription_state should be true so subscribers can progress independently")
	}
	if !plan.SupportsFiltering {
		issues = append(issues, "supports_filtering should be true when not every subscriber needs every event")
	}
	if !plan.HasReplayCursor {
		issues = append(issues, "has_replay_cursor should be true so one subscriber can recover without global impact")
	}
	if !plan.HasIsolationControls {
		issues = append(issues, "has_isolation_controls should be true so slow subscribers do not harm fast ones")
	}
	if !plan.HasBacklogBudgets {
		issues = append(issues, "has_backlog_budgets should be true so storage and delivery lag stay bounded")
	}
	if !plan.HasDeadLetterPath {
		issues = append(issues, "has_dead_letter_path should be true for repeatedly failing deliveries")
	}
	return issues
}

func main() {
	name := flag.String("name", "pubsub-fanout", "plan name")
	flag.Parse()

	plan := PubSubPlan{
		Name:                    *name,
		HasTopicPartitions:      true,
		HasPerSubscriptionState: true,
		SupportsFiltering:       true,
		HasReplayCursor:         true,
		HasIsolationControls:    true,
		HasBacklogBudgets:       true,
		HasDeadLetterPath:       true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidatePubSubPlan(plan),
	})
}
