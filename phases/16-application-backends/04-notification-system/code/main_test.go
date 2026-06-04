package main

import "testing"

func TestValidateNotificationPlanHealthy(t *testing.T) {
	plan := NotificationPlan{
		Name:                     "healthy",
		HasDedupeStore:           true,
		HasPreferenceChecks:      true,
		UsesPriorityQueues:       true,
		HasRetryBudget:           true,
		SupportsProviderFailover: true,
		QuietHoursEnforced:       true,
		TracksDeliveryState:      true,
	}
	if issues := ValidateNotificationPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateNotificationPlan returned issues: %v", issues)
	}
}

func TestValidateNotificationPlanRejectsWeakPlan(t *testing.T) {
	plan := NotificationPlan{Name: "weak"}
	if issues := ValidateNotificationPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateNotificationPlan returned too few issues: %v", issues)
	}
}
