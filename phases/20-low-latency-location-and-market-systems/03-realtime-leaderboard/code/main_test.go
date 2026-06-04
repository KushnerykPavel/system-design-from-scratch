package main

import "testing"

func TestValidateLeaderboardConfigHealthy(t *testing.T) {
	cfg := LeaderboardConfig{
		Name:                   "healthy",
		ShardCount:             32,
		TopNCacheSize:          100,
		AroundMeWindow:         25,
		PublishDelaySeconds:    5,
		ValidationQueueEnabled: true,
		CorrectionSupport:      true,
	}
	if issues := ValidateLeaderboardConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateLeaderboardConfig returned issues: %v", issues)
	}
}

func TestValidateLeaderboardConfigWeak(t *testing.T) {
	cfg := LeaderboardConfig{
		Name:                "weak",
		ShardCount:          1,
		TopNCacheSize:       1,
		AroundMeWindow:      500,
		PublishDelaySeconds: 1000,
	}
	if issues := ValidateLeaderboardConfig(cfg); len(issues) < 5 {
		t.Fatalf("ValidateLeaderboardConfig returned too few issues: %v", issues)
	}
}
