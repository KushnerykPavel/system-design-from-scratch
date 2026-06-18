package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)

// ItemEmbedding stores a precomputed item vector.
type ItemEmbedding struct {
	ItemID    string
	Embedding []float64
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Ranked is a candidate item with its score.
type Ranked struct {
	ItemID string  `json:"item_id"`
	Score  float64 `json:"score"`
}

// Retrieve returns the top-K items most similar to the user embedding.
func Retrieve(userEmbedding []float64, items []ItemEmbedding, topK int) []Ranked {
	ranked := make([]Ranked, 0, len(items))
	for _, item := range items {
		score := cosineSimilarity(userEmbedding, item.Embedding)
		ranked = append(ranked, Ranked{ItemID: item.ItemID, Score: score})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})
	if topK > len(ranked) {
		topK = len(ranked)
	}
	return ranked[:topK]
}

func main() {
	// Simulate a 4-dimensional embedding space.
	items := []ItemEmbedding{
		{ItemID: "action-1", Embedding: []float64{0.9, 0.1, 0.0, 0.2}},
		{ItemID: "drama-1", Embedding: []float64{0.1, 0.9, 0.1, 0.0}},
		{ItemID: "action-2", Embedding: []float64{0.8, 0.2, 0.1, 0.3}},
		{ItemID: "comedy-1", Embedding: []float64{0.0, 0.1, 0.9, 0.1}},
		{ItemID: "action-3", Embedding: []float64{0.85, 0.05, 0.0, 0.1}},
	}
	// User embedding is action-oriented.
	userEmbedding := []float64{0.95, 0.05, 0.0, 0.1}
	top3 := Retrieve(userEmbedding, items, 3)

	result := map[string]any{
		"user_embedding": userEmbedding,
		"top_3":          top3,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
