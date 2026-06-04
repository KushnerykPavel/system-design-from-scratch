package main

import "testing"

func TestValidateCrawlPlanHealthy(t *testing.T) {
	plan := CrawlPlan{
		Name:                "healthy",
		Fetchers:            100,
		PerHostQPS:          1,
		QueueLeaseSeconds:   90,
		RobotsEnforced:      true,
		FrontierReplicas:    3,
		ContentDedupEnabled: true,
		RecrawlTiers:        2,
	}
	if issues := ValidateCrawlPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateCrawlPlan returned issues: %v", issues)
	}
}

func TestValidateCrawlPlanWeak(t *testing.T) {
	plan := CrawlPlan{
		Name:              "weak",
		Fetchers:          1,
		PerHostQPS:        0,
		QueueLeaseSeconds: 10,
		FrontierReplicas:  1,
		RecrawlTiers:      1,
	}
	if issues := ValidateCrawlPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateCrawlPlan returned too few issues: %v", issues)
	}
}
