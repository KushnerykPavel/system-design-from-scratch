package main

import "testing"

func TestValidatePubSubPlan(t *testing.T) {
	full := PubSubPlan{
		HasTopicPartitions:      true,
		HasPerSubscriptionState: true,
		SupportsFiltering:       true,
		HasReplayCursor:         true,
		HasIsolationControls:    true,
		HasBacklogBudgets:       true,
		HasDeadLetterPath:       true,
	}
	if issues := ValidatePubSubPlan(full); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	full.HasIsolationControls = false
	if issues := ValidatePubSubPlan(full); len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
}
