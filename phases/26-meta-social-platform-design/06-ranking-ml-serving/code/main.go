package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FeatureWeights controls how much each feature dimension contributes to the final score.
type FeatureWeights struct {
	Affinity   float64 `json:"affinity"`
	Recency    float64 `json:"recency"`
	Engagement float64 `json:"engagement"`
}

// FeedItem represents a candidate post entering the ranking pipeline.
type FeedItem struct {
	ItemID         string  `json:"item_id"`
	AffinityScore  float64 `json:"affinity_score"`  // user-item affinity (0.0–1.0)
	RecencyScore   float64 `json:"recency_score"`   // how recent the post is (0.0–1.0)
	EngagementRate float64 `json:"engagement_rate"` // normalized engagement (0.0–1.0)
}

// rankedItem is an internal type that attaches the computed score to a FeedItem.
type rankedItem struct {
	FeedItem
	Score float64 `json:"score"`
}

// RankItems sorts a slice of FeedItems by weighted composite score descending.
// It returns a new slice — the original is not modified.
func RankItems(items []FeedItem, weights FeatureWeights) []FeedItem {
	scored := make([]rankedItem, len(items))
	for i, item := range items {
		score := weights.Affinity*item.AffinityScore +
			weights.Recency*item.RecencyScore +
			weights.Engagement*item.EngagementRate
		scored[i] = rankedItem{FeedItem: item, Score: score}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		// Tie-break: prefer higher affinity, then lower item ID lexicographically.
		if scored[i].AffinityScore != scored[j].AffinityScore {
			return scored[i].AffinityScore > scored[j].AffinityScore
		}
		return scored[i].ItemID < scored[j].ItemID
	})

	result := make([]FeedItem, len(scored))
	for i, s := range scored {
		result[i] = s.FeedItem
	}
	return result
}

// computeScore returns the weighted score for a single item — useful for display and tests.
func computeScore(item FeedItem, weights FeatureWeights) float64 {
	return weights.Affinity*item.AffinityScore +
		weights.Recency*item.RecencyScore +
		weights.Engagement*item.EngagementRate
}

func main() {
	candidates := []FeedItem{
		{ItemID: "post-001", AffinityScore: 0.9, RecencyScore: 0.8, EngagementRate: 0.3},
		{ItemID: "post-002", AffinityScore: 0.5, RecencyScore: 0.95, EngagementRate: 0.7},
		{ItemID: "post-003", AffinityScore: 0.2, RecencyScore: 0.6, EngagementRate: 0.9},
		{ItemID: "post-004", AffinityScore: 0.8, RecencyScore: 0.3, EngagementRate: 0.5},
		{ItemID: "post-005", AffinityScore: 0.6, RecencyScore: 0.7, EngagementRate: 0.6},
	}

	// Profile A: personalization-heavy (affinity dominates).
	profileA := FeatureWeights{Affinity: 0.6, Recency: 0.2, Engagement: 0.2}

	// Profile B: engagement-driven (good for cold start / viral content discovery).
	profileB := FeatureWeights{Affinity: 0.2, Recency: 0.2, Engagement: 0.6}

	rankedA := RankItems(candidates, profileA)
	rankedB := RankItems(candidates, profileB)

	type rankedWithScore struct {
		ItemID string  `json:"item_id"`
		Score  float64 `json:"score"`
	}

	toScored := func(items []FeedItem, w FeatureWeights) []rankedWithScore {
		out := make([]rankedWithScore, len(items))
		for i, it := range items {
			out[i] = rankedWithScore{ItemID: it.ItemID, Score: computeScore(it, w)}
		}
		return out
	}

	output := map[string]any{
		"profile_a_personalization": map[string]any{
			"weights": profileA,
			"ranked":  toScored(rankedA, profileA),
		},
		"profile_b_engagement_driven": map[string]any{
			"weights": profileB,
			"ranked":  toScored(rankedB, profileB),
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
