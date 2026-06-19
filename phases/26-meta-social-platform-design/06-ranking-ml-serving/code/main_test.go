package main

import (
	"testing"
)

func TestRankItemsBasicOrder(t *testing.T) {
	items := []FeedItem{
		{ItemID: "low", AffinityScore: 0.1, RecencyScore: 0.1, EngagementRate: 0.1},
		{ItemID: "high", AffinityScore: 0.9, RecencyScore: 0.9, EngagementRate: 0.9},
		{ItemID: "mid", AffinityScore: 0.5, RecencyScore: 0.5, EngagementRate: 0.5},
	}
	weights := FeatureWeights{Affinity: 1.0 / 3, Recency: 1.0 / 3, Engagement: 1.0 / 3}
	ranked := RankItems(items, weights)

	if ranked[0].ItemID != "high" {
		t.Errorf("expected 'high' first, got %s", ranked[0].ItemID)
	}
	if ranked[1].ItemID != "mid" {
		t.Errorf("expected 'mid' second, got %s", ranked[1].ItemID)
	}
	if ranked[2].ItemID != "low" {
		t.Errorf("expected 'low' third, got %s", ranked[2].ItemID)
	}
}

func TestRankItemsAffinityWeightDominates(t *testing.T) {
	// post-a has high affinity but low recency and engagement.
	// post-b has low affinity but high recency and engagement.
	// With affinity weight = 0.9, post-a should win.
	items := []FeedItem{
		{ItemID: "post-a", AffinityScore: 0.9, RecencyScore: 0.1, EngagementRate: 0.1},
		{ItemID: "post-b", AffinityScore: 0.1, RecencyScore: 0.9, EngagementRate: 0.9},
	}
	weights := FeatureWeights{Affinity: 0.9, Recency: 0.05, Engagement: 0.05}
	ranked := RankItems(items, weights)

	if ranked[0].ItemID != "post-a" {
		t.Errorf("with high affinity weight, expected post-a first, got %s", ranked[0].ItemID)
	}
}

func TestRankItemsEngagementWeightDominates(t *testing.T) {
	// With engagement weight = 0.9, post-b (high engagement) should win.
	items := []FeedItem{
		{ItemID: "post-a", AffinityScore: 0.9, RecencyScore: 0.1, EngagementRate: 0.1},
		{ItemID: "post-b", AffinityScore: 0.1, RecencyScore: 0.9, EngagementRate: 0.9},
	}
	weights := FeatureWeights{Affinity: 0.05, Recency: 0.05, Engagement: 0.9}
	ranked := RankItems(items, weights)

	if ranked[0].ItemID != "post-b" {
		t.Errorf("with high engagement weight, expected post-b first, got %s", ranked[0].ItemID)
	}
}

func TestRankItemsPreservesAllItems(t *testing.T) {
	items := []FeedItem{
		{ItemID: "a", AffinityScore: 0.5, RecencyScore: 0.5, EngagementRate: 0.5},
		{ItemID: "b", AffinityScore: 0.3, RecencyScore: 0.7, EngagementRate: 0.2},
		{ItemID: "c", AffinityScore: 0.8, RecencyScore: 0.1, EngagementRate: 0.9},
		{ItemID: "d", AffinityScore: 0.2, RecencyScore: 0.2, EngagementRate: 0.2},
	}
	weights := FeatureWeights{Affinity: 0.4, Recency: 0.3, Engagement: 0.3}
	ranked := RankItems(items, weights)

	if len(ranked) != len(items) {
		t.Fatalf("expected %d items in result, got %d", len(items), len(ranked))
	}
}

func TestRankItemsDoesNotMutateOriginal(t *testing.T) {
	items := []FeedItem{
		{ItemID: "first", AffinityScore: 0.1, RecencyScore: 0.1, EngagementRate: 0.1},
		{ItemID: "second", AffinityScore: 0.9, RecencyScore: 0.9, EngagementRate: 0.9},
	}
	original := items[0].ItemID
	weights := FeatureWeights{Affinity: 0.4, Recency: 0.3, Engagement: 0.3}
	RankItems(items, weights)

	if items[0].ItemID != original {
		t.Errorf("original slice was mutated: expected %s at index 0, got %s", original, items[0].ItemID)
	}
}

func TestRankItemsTieBreaking(t *testing.T) {
	// Two items with identical scores — tie-break by affinity, then by item ID.
	items := []FeedItem{
		{ItemID: "z-post", AffinityScore: 0.5, RecencyScore: 0.5, EngagementRate: 0.5},
		{ItemID: "a-post", AffinityScore: 0.5, RecencyScore: 0.5, EngagementRate: 0.5},
	}
	weights := FeatureWeights{Affinity: 1.0 / 3, Recency: 1.0 / 3, Engagement: 1.0 / 3}
	ranked := RankItems(items, weights)

	// Scores are identical, affinities are identical, so tie-break goes to lexicographically smaller ID.
	if ranked[0].ItemID != "a-post" {
		t.Errorf("expected 'a-post' first on tie-break, got %s", ranked[0].ItemID)
	}
}

func TestRankItemsEmptySlice(t *testing.T) {
	weights := FeatureWeights{Affinity: 0.4, Recency: 0.3, Engagement: 0.3}
	ranked := RankItems([]FeedItem{}, weights)
	if len(ranked) != 0 {
		t.Fatalf("expected empty result for empty input, got %d items", len(ranked))
	}
}

func TestComputeScore(t *testing.T) {
	item := FeedItem{ItemID: "x", AffinityScore: 0.8, RecencyScore: 0.5, EngagementRate: 0.4}
	weights := FeatureWeights{Affinity: 0.5, Recency: 0.3, Engagement: 0.2}
	// expected = 0.5*0.8 + 0.3*0.5 + 0.2*0.4 = 0.40 + 0.15 + 0.08 = 0.63
	expected := 0.63
	got := computeScore(item, weights)
	if abs(got-expected) > 1e-9 {
		t.Errorf("computeScore = %f, want %f", got, expected)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
