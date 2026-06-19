package main

import "testing"

const testNow int64 = 1700000000

func TestFilterSpamRemovesHighScoreItems(t *testing.T) {
	items := []FeedItem{
		{ID: "clean", SpamScore: 0.1},
		{ID: "spam", SpamScore: 0.8},
		{ID: "borderline", SpamScore: 0.5}, // at threshold, not removed (< threshold required)
	}
	filtered := FilterSpam(items, 0.5)
	for _, it := range filtered {
		if it.ID == "spam" {
			t.Fatal("spam item with SpamScore=0.8 should have been filtered out")
		}
	}
	// borderline (0.5) should be filtered since 0.5 is not < 0.5
	for _, it := range filtered {
		if it.ID == "borderline" {
			t.Fatal("borderline item with SpamScore=0.5 should be filtered (not < threshold)")
		}
	}
	if len(filtered) != 1 {
		t.Fatalf("expected 1 item after filtering, got %d", len(filtered))
	}
}

func TestFilterSpamKeepsCleanItems(t *testing.T) {
	items := []FeedItem{
		{ID: "a", SpamScore: 0.0},
		{ID: "b", SpamScore: 0.2},
		{ID: "c", SpamScore: 0.49},
	}
	filtered := FilterSpam(items, 0.5)
	if len(filtered) != 3 {
		t.Fatalf("expected all 3 clean items to pass filter, got %d", len(filtered))
	}
}

func TestRankFeedOrderHighEngagementAndRecentFirst(t *testing.T) {
	weights := DefaultWeights()
	items := []FeedItem{
		{
			ID: "old-low", CreatedAtUnix: testNow - 80000, // very old
			EngagementRate: 0.1, DegreeProximity: 0.2, SpamScore: 0.0, ViralCoefficient: 0.0,
		},
		{
			ID: "new-high", CreatedAtUnix: testNow - 600, // 10 min ago
			EngagementRate: 0.9, DegreeProximity: 1.0, SpamScore: 0.0, ViralCoefficient: 0.0,
		},
	}
	ranked := RankFeed(items, weights, testNow)
	if ranked[0].ID != "new-high" {
		t.Fatalf("expected new-high to rank first, got %s", ranked[0].ID)
	}
	if ranked[1].ID != "old-low" {
		t.Fatalf("expected old-low to rank last, got %s", ranked[1].ID)
	}
}

func TestRankFeedAntiViralityPenalty(t *testing.T) {
	weights := DefaultWeights()
	// viral has high engagement but viral coefficient triggers penalty
	// safe has moderate engagement but no penalty
	items := []FeedItem{
		{
			ID: "viral", CreatedAtUnix: testNow - 3600,
			EngagementRate: 0.9, DegreeProximity: 0.5, SpamScore: 0.0,
			ViralCoefficient: 3.5, // triggers -0.2 penalty
		},
		{
			ID: "safe", CreatedAtUnix: testNow - 3600,
			EngagementRate: 0.7, DegreeProximity: 0.5, SpamScore: 0.0,
			ViralCoefficient: 1.0, // no penalty
		},
	}
	ranked := RankFeed(items, weights, testNow)
	// safe should rank above viral because viral penalty offsets engagement advantage
	viralScore := ScoreItem(items[0], weights, testNow)
	safeScore := ScoreItem(items[1], weights, testNow)
	if viralScore >= safeScore {
		t.Fatalf("expected viral (score=%.4f) to score lower than safe (score=%.4f) due to anti-virality penalty", viralScore, safeScore)
	}
	if ranked[0].ID != "safe" {
		t.Fatalf("expected safe to rank first after anti-virality penalty, got %s", ranked[0].ID)
	}
}

func TestRankFeedProximityImpact(t *testing.T) {
	weights := DefaultWeights()
	// high-proximity post (1st-degree author) vs low-proximity post (3rd-degree)
	items := []FeedItem{
		{
			ID: "third-degree", CreatedAtUnix: testNow - 3600,
			EngagementRate: 0.8, DegreeProximity: 0.2, SpamScore: 0.0, ViralCoefficient: 0.0,
		},
		{
			ID: "first-degree", CreatedAtUnix: testNow - 3600,
			EngagementRate: 0.5, DegreeProximity: 1.0, SpamScore: 0.0, ViralCoefficient: 0.0,
		},
	}
	ranked := RankFeed(items, weights, testNow)
	// first-degree should win because proximity weight (0.45) × 1.0 > third-degree proximity 0.45 × 0.2
	// first-degree score: 0.2*recency + 0.35*0.5 + 0.45*1.0 = recency+0.175+0.45
	// third-degree score: 0.2*recency + 0.35*0.8 + 0.45*0.2 = recency+0.28+0.09
	// first-degree > third-degree because 0.625 > 0.37 (ignoring equal recency)
	if ranked[0].ID != "first-degree" {
		t.Fatalf("expected first-degree to rank above third-degree, got %s", ranked[0].ID)
	}
}

func TestScoreItemRecencyDecay(t *testing.T) {
	weights := DefaultWeights()
	fresh := FeedItem{
		ID: "fresh", CreatedAtUnix: testNow - 1800, // 30 min ago
		EngagementRate: 0.5, DegreeProximity: 0.5, SpamScore: 0.0, ViralCoefficient: 0.0,
	}
	stale := FeedItem{
		ID: "stale", CreatedAtUnix: testNow - 86400, // 24h ago
		EngagementRate: 0.5, DegreeProximity: 0.5, SpamScore: 0.0, ViralCoefficient: 0.0,
	}
	freshScore := ScoreItem(fresh, weights, testNow)
	staleScore := ScoreItem(stale, weights, testNow)
	if freshScore <= staleScore {
		t.Fatalf("expected fresh item (score=%.4f) to score higher than stale (score=%.4f)", freshScore, staleScore)
	}
}
