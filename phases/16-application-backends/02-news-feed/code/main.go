package main

import (
	"encoding/json"
	"flag"
	"os"
)

type FeedPlan struct {
	Name                    string `json:"name"`
	CelebrityThreshold      int    `json:"celebrity_threshold"`
	UsesPushFanout          bool   `json:"uses_push_fanout"`
	UsesPullFanout          bool   `json:"uses_pull_fanout"`
	HasRankingFallback      bool   `json:"has_ranking_fallback"`
	ModerationTombstones    bool   `json:"moderation_tombstones"`
	FreshnessTargetSeconds  int    `json:"freshness_target_seconds"`
	CoalescesTimelineMisses bool   `json:"coalesces_timeline_misses"`
}

func ValidateFeedPlan(plan FeedPlan) []string {
	var issues []string
	if plan.CelebrityThreshold <= 0 {
		issues = append(issues, "celebrity_threshold must be positive")
	}
	if !plan.UsesPushFanout || !plan.UsesPullFanout {
		issues = append(issues, "uses_push_fanout and uses_pull_fanout should both be true for mixed-skew timelines")
	}
	if !plan.HasRankingFallback {
		issues = append(issues, "has_ranking_fallback should be true so feed reads survive ranking degradation")
	}
	if !plan.ModerationTombstones {
		issues = append(issues, "moderation_tombstones should be enabled for fast hide and delete propagation")
	}
	if plan.FreshnessTargetSeconds <= 0 || plan.FreshnessTargetSeconds > 300 {
		issues = append(issues, "freshness_target_seconds should be a realistic positive bound")
	}
	if !plan.CoalescesTimelineMisses {
		issues = append(issues, "coalesces_timeline_misses should be true to limit cache stampedes")
	}
	return issues
}

func main() {
	name := flag.String("name", "hybrid-feed", "plan name")
	flag.Parse()

	plan := FeedPlan{
		Name:                    *name,
		CelebrityThreshold:      100000,
		UsesPushFanout:          true,
		UsesPullFanout:          true,
		HasRankingFallback:      true,
		ModerationTombstones:    true,
		FreshnessTargetSeconds:  10,
		CoalescesTimelineMisses: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateFeedPlan(plan),
	})
}
