package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FeedItem represents a candidate post entering the LinkedIn feed ranking pipeline.
type FeedItem struct {
	ID               string  `json:"id"`
	AuthorID         string  `json:"author_id"`
	CreatedAtUnix    int64   `json:"created_at_unix"`
	EngagementRate   float64 `json:"engagement_rate"`    // normalized 0.0–1.0
	DegreeProximity  float64 `json:"degree_proximity"`   // 1.0 = 1st-degree, 0.5 = 2nd-degree, 0.2 = 3rd-degree
	SpamScore        float64 `json:"spam_score"`         // 0.0 = clean, 1.0 = spam
	ViralCoefficient float64 `json:"viral_coefficient"`  // reshare rate; >2.5 triggers penalty
	ContentType      string  `json:"content_type"`       // "post", "article", "job_alert", "milestone"
}

// RankingWeights controls the contribution of each signal to the final score.
type RankingWeights struct {
	Recency         float64 `json:"recency"`
	Engagement      float64 `json:"engagement"`
	Proximity       float64 `json:"proximity"`
	QualityPenalty  float64 `json:"quality_penalty"`   // multiplier applied to viral coefficient penalty
}

// DefaultWeights returns LinkedIn-tuned ranking weights.
// Professional proximity and engagement are balanced;
// recency gets a moderate weight to surface fresh content.
func DefaultWeights() RankingWeights {
	return RankingWeights{
		Recency:        0.2,
		Engagement:     0.35,
		Proximity:      0.45,
		QualityPenalty: 0.2,
	}
}

// ScoreItem computes the ranking score for a single feed item.
// Recency score decays linearly over 24 hours from the reference time.
// Viral coefficient above 2.5 applies a penalty proportional to QualityPenalty weight.
func ScoreItem(item FeedItem, weights RankingWeights, nowUnix int64) float64 {
	// Recency: 1.0 if just posted, 0.0 if 24h+ old.
	ageHours := float64(nowUnix-item.CreatedAtUnix) / 3600.0
	recency := 1.0 - (ageHours / 24.0)
	if recency < 0 {
		recency = 0
	}

	score := weights.Recency*recency +
		weights.Engagement*item.EngagementRate +
		weights.Proximity*item.DegreeProximity

	// Anti-virality penalty: posts spreading too fast get deprioritized.
	if item.ViralCoefficient > 2.5 {
		score -= weights.QualityPenalty
	}

	return score
}

// FilterSpam removes items with spam score at or above the given threshold.
func FilterSpam(items []FeedItem, threshold float64) []FeedItem {
	out := make([]FeedItem, 0, len(items))
	for _, item := range items {
		if item.SpamScore < threshold {
			out = append(out, item)
		}
	}
	return out
}

// RankFeed scores and sorts feed items in descending score order.
// Returns a new slice; the original is not modified.
func RankFeed(items []FeedItem, weights RankingWeights, nowUnix int64) []FeedItem {
	type scored struct {
		item  FeedItem
		score float64
	}
	ss := make([]scored, len(items))
	for i, it := range items {
		ss[i] = scored{item: it, score: ScoreItem(it, weights, nowUnix)}
	}
	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].score > ss[j].score
	})
	result := make([]FeedItem, len(ss))
	for i, s := range ss {
		result[i] = s.item
	}
	return result
}

func main() {
	// Reference time: treat all items relative to "now = 1700000000".
	const now int64 = 1700000000

	items := []FeedItem{
		{
			ID: "post-001", AuthorID: "alice",
			CreatedAtUnix: now - 3600, // 1 hour ago
			EngagementRate: 0.8, DegreeProximity: 1.0,
			SpamScore: 0.05, ViralCoefficient: 0.8,
			ContentType: "post",
		},
		{
			ID: "post-002", AuthorID: "bob",
			CreatedAtUnix: now - 7200, // 2 hours ago
			EngagementRate: 0.3, DegreeProximity: 0.5,
			SpamScore: 0.1, ViralCoefficient: 0.5,
			ContentType: "article",
		},
		{
			ID: "post-003", AuthorID: "carol",
			CreatedAtUnix: now - 1800, // 30 minutes ago
			EngagementRate: 0.9, DegreeProximity: 0.5,
			SpamScore: 0.05, ViralCoefficient: 3.2, // triggers anti-virality penalty
			ContentType: "post",
		},
		{
			ID: "post-004", AuthorID: "dave",
			CreatedAtUnix: now - 86400, // 24 hours ago — old
			EngagementRate: 0.2, DegreeProximity: 0.2,
			SpamScore: 0.6, ViralCoefficient: 0.1, // spam: filtered out
			ContentType: "post",
		},
		{
			ID: "post-005", AuthorID: "eve",
			CreatedAtUnix: now - 600, // 10 minutes ago
			EngagementRate: 0.6, DegreeProximity: 1.0,
			SpamScore: 0.02, ViralCoefficient: 1.2,
			ContentType: "milestone",
		},
	}

	weights := DefaultWeights()
	spamThreshold := 0.5

	filtered := FilterSpam(items, spamThreshold)
	ranked := RankFeed(filtered, weights, now)

	type rankedOutput struct {
		ID              string  `json:"id"`
		AuthorID        string  `json:"author_id"`
		Score           float64 `json:"score"`
		SpamScore       float64 `json:"spam_score"`
		ViralCoefficient float64 `json:"viral_coefficient"`
	}

	out := make([]rankedOutput, len(ranked))
	for i, it := range ranked {
		out[i] = rankedOutput{
			ID:              it.ID,
			AuthorID:        it.AuthorID,
			Score:           ScoreItem(it, weights, now),
			SpamScore:       it.SpamScore,
			ViralCoefficient: it.ViralCoefficient,
		}
	}

	output := map[string]any{
		"spam_threshold":     spamThreshold,
		"items_before_filter": len(items),
		"items_after_filter":  len(filtered),
		"ranked_feed":         out,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
