package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Item represents a DynamoDB item with a partition key, sort key, and entity type.
type Item struct {
	PK         string
	SK         string
	EntityType string
}

// PartitionStats holds the result of analyzing partition key distribution.
// DistributionScore is 0–100 where 100 means perfectly uniform distribution.
// Formula: score = 100 * max(0, 1 - stddev/mean), clamped to [0, 100].
// A score below 70 indicates hot-partition risk.
type PartitionStats struct {
	PartitionCounts   map[string]int `json:"partition_counts"`
	TotalItems        int            `json:"total_items"`
	HotPartitions     []string       `json:"hot_partitions"`
	DistributionScore int            `json:"distribution_score"`
}

// AnalyzeDistribution computes partition key distribution statistics.
// A partition is flagged as hot if it holds more than 10% of total items.
// DistributionScore = 100 * max(0, 1 - stddev/mean), clamped to [0,100].
func AnalyzeDistribution(items []Item) PartitionStats {
	if len(items) == 0 {
		return PartitionStats{
			PartitionCounts:   map[string]int{},
			DistributionScore: 100,
		}
	}

	counts := make(map[string]int)
	for _, it := range items {
		counts[it.PK]++
	}

	total := len(items)
	n := float64(len(counts))

	// Compute mean and stddev of counts across partition keys.
	var sum float64
	for _, c := range counts {
		sum += float64(c)
	}
	mean := sum / n

	var variance float64
	for _, c := range counts {
		diff := float64(c) - mean
		variance += diff * diff
	}
	variance /= n
	stddev := math.Sqrt(variance)

	// DistributionScore: 100 means perfectly uniform (stddev=0).
	score := 100.0
	if mean > 0 {
		score = 100.0 * math.Max(0, 1.0-stddev/mean)
	}
	if score > 100 {
		score = 100
	}

	// Identify hot partitions: any PK that holds > 10% of total items.
	threshold := float64(total) * 0.10
	var hot []string
	for pk, c := range counts {
		if float64(c) > threshold {
			hot = append(hot, fmt.Sprintf("%s (%d items, %.1f%%)", pk, c, float64(c)*100/float64(total)))
		}
	}

	return PartitionStats{
		PartitionCounts:   counts,
		TotalItems:        total,
		HotPartitions:     hot,
		DistributionScore: int(math.Round(score)),
	}
}

func main() {
	// --- Sample 1: Uniform distribution (10 partitions, 10 items each) ---
	uniform := make([]Item, 0, 100)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			uniform = append(uniform, Item{
				PK:         fmt.Sprintf("PRODUCT#pk-%02d", i),
				SK:         fmt.Sprintf("ITEM#%03d", j),
				EntityType: "product",
			})
		}
	}

	// --- Sample 2: Hot partition — one key holds 60% of traffic ---
	skewed := make([]Item, 0, 100)
	for j := 0; j < 60; j++ {
		skewed = append(skewed, Item{PK: "PRODUCT#viral-item", SK: fmt.Sprintf("VIEW#%03d", j), EntityType: "product"})
	}
	for i := 1; i < 5; i++ {
		for j := 0; j < 10; j++ {
			skewed = append(skewed, Item{
				PK:         fmt.Sprintf("PRODUCT#pk-%02d", i),
				SK:         fmt.Sprintf("VIEW#%03d", j),
				EntityType: "product",
			})
		}
	}

	// --- Sample 3: Bad sequential keys (order ID prefix) ---
	sequential := make([]Item, 0, 50)
	for i := 0; i < 50; i++ {
		sequential = append(sequential, Item{
			PK:         fmt.Sprintf("ORDER#%05d", i),
			SK:         "METADATA",
			EntityType: "order",
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	fmt.Println("=== Uniform distribution (expected score ~100) ===")
	_ = enc.Encode(AnalyzeDistribution(uniform))

	fmt.Println("=== Skewed distribution (hot partition expected) ===")
	_ = enc.Encode(AnalyzeDistribution(skewed))

	fmt.Println("=== Sequential keys (each unique, no hot partition, score ~100) ===")
	_ = enc.Encode(AnalyzeDistribution(sequential))
}
