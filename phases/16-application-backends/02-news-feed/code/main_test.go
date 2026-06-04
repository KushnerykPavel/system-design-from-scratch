package main

import "testing"

func TestValidateFeedPlanHealthy(t *testing.T) {
	plan := FeedPlan{
		Name:                    "healthy",
		CelebrityThreshold:      1000,
		UsesPushFanout:          true,
		UsesPullFanout:          true,
		HasRankingFallback:      true,
		ModerationTombstones:    true,
		FreshnessTargetSeconds:  15,
		CoalescesTimelineMisses: true,
	}
	if issues := ValidateFeedPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateFeedPlan returned issues: %v", issues)
	}
}

func TestValidateFeedPlanRejectsWeakPlan(t *testing.T) {
	plan := FeedPlan{
		Name:                   "weak",
		CelebrityThreshold:     0,
		UsesPushFanout:         true,
		FreshnessTargetSeconds: 600,
	}
	if issues := ValidateFeedPlan(plan); len(issues) < 4 {
		t.Fatalf("ValidateFeedPlan returned too few issues: %v", issues)
	}
}
