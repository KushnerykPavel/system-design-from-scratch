package main

import "testing"

func TestValidateMarketDataPolicyHealthy(t *testing.T) {
	cfg := MarketDataPolicy{
		Name:                    "healthy",
		PremiumReplayMinutes:    15,
		StandardReplayMinutes:   1,
		MaxSymbolsPerSubscriber: 500,
		RegionalRelays:          3,
		ReplayLaneIsolation:     true,
		SubscriberQuotaEnabled:  true,
	}
	if issues := ValidateMarketDataPolicy(cfg); len(issues) != 0 {
		t.Fatalf("ValidateMarketDataPolicy returned issues: %v", issues)
	}
}

func TestValidateMarketDataPolicyWeak(t *testing.T) {
	cfg := MarketDataPolicy{
		Name:                    "weak",
		PremiumReplayMinutes:    0,
		StandardReplayMinutes:   5,
		MaxSymbolsPerSubscriber: 1,
		RegionalRelays:          1,
	}
	if issues := ValidateMarketDataPolicy(cfg); len(issues) < 5 {
		t.Fatalf("ValidateMarketDataPolicy returned too few issues: %v", issues)
	}
}
