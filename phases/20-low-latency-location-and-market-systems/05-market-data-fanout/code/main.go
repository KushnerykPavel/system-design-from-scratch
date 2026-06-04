package main

import (
	"encoding/json"
	"flag"
	"os"
)

type MarketDataPolicy struct {
	Name                    string `json:"name"`
	PremiumReplayMinutes    int    `json:"premium_replay_minutes"`
	StandardReplayMinutes   int    `json:"standard_replay_minutes"`
	MaxSymbolsPerSubscriber int    `json:"max_symbols_per_subscriber"`
	RegionalRelays          int    `json:"regional_relays"`
	ReplayLaneIsolation     bool   `json:"replay_lane_isolation"`
	SubscriberQuotaEnabled  bool   `json:"subscriber_quota_enabled"`
}

func ValidateMarketDataPolicy(cfg MarketDataPolicy) []string {
	var issues []string
	if cfg.PremiumReplayMinutes < cfg.StandardReplayMinutes {
		issues = append(issues, "premium_replay_minutes should be at least standard_replay_minutes")
	}
	if cfg.StandardReplayMinutes < 0 || cfg.StandardReplayMinutes > 60 {
		issues = append(issues, "standard_replay_minutes should stay between 0 and 60")
	}
	if cfg.MaxSymbolsPerSubscriber < 10 {
		issues = append(issues, "max_symbols_per_subscriber should be at least 10")
	}
	if cfg.RegionalRelays < 2 {
		issues = append(issues, "regional_relays should be at least 2")
	}
	if !cfg.ReplayLaneIsolation {
		issues = append(issues, "replay_lane_isolation should be enabled")
	}
	if !cfg.SubscriberQuotaEnabled {
		issues = append(issues, "subscriber_quota_enabled should be enabled")
	}
	return issues
}

func main() {
	name := flag.String("name", "market-data", "policy name")
	flag.Parse()

	cfg := MarketDataPolicy{
		Name:                    *name,
		PremiumReplayMinutes:    15,
		StandardReplayMinutes:   1,
		MaxSymbolsPerSubscriber: 500,
		RegionalRelays:          3,
		ReplayLaneIsolation:     true,
		SubscriberQuotaEnabled:  true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateMarketDataPolicy(cfg),
	})
}
