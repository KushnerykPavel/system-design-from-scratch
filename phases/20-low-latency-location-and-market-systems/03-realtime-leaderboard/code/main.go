package main

import (
	"encoding/json"
	"flag"
	"os"
)

type LeaderboardConfig struct {
	Name                   string `json:"name"`
	ShardCount             int    `json:"shard_count"`
	TopNCacheSize          int    `json:"top_n_cache_size"`
	AroundMeWindow         int    `json:"around_me_window"`
	PublishDelaySeconds    int    `json:"publish_delay_seconds"`
	ValidationQueueEnabled bool   `json:"validation_queue_enabled"`
	CorrectionSupport      bool   `json:"correction_support"`
}

func ValidateLeaderboardConfig(cfg LeaderboardConfig) []string {
	var issues []string
	if cfg.ShardCount < 8 {
		issues = append(issues, "shard_count should be at least 8")
	}
	if cfg.TopNCacheSize < 10 || cfg.TopNCacheSize > 1000 {
		issues = append(issues, "top_n_cache_size should stay between 10 and 1000")
	}
	if cfg.AroundMeWindow < 5 || cfg.AroundMeWindow > 200 {
		issues = append(issues, "around_me_window should stay between 5 and 200")
	}
	if cfg.PublishDelaySeconds < 0 || cfg.PublishDelaySeconds > 300 {
		issues = append(issues, "publish_delay_seconds should stay between 0 and 300")
	}
	if !cfg.ValidationQueueEnabled {
		issues = append(issues, "validation_queue_enabled should usually be true for suspicious updates")
	}
	if !cfg.CorrectionSupport {
		issues = append(issues, "correction_support should be enabled for reversals or moderation")
	}
	return issues
}

func main() {
	name := flag.String("name", "realtime-leaderboard", "config name")
	flag.Parse()

	cfg := LeaderboardConfig{
		Name:                   *name,
		ShardCount:             32,
		TopNCacheSize:          100,
		AroundMeWindow:         25,
		PublishDelaySeconds:    5,
		ValidationQueueEnabled: true,
		CorrectionSupport:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateLeaderboardConfig(cfg),
	})
}
