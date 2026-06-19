package main

import (
	"testing"
)

// TestUniformDistributionScore verifies that perfectly uniform partition distribution
// yields a DistributionScore of 100.
func TestUniformDistributionScore(t *testing.T) {
	items := make([]Item, 0, 100)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			items = append(items, Item{
				PK:         "PK#" + string(rune('A'+i)),
				SK:         "SK#" + string(rune('a'+j)),
				EntityType: "test",
			})
		}
	}
	stats := AnalyzeDistribution(items)
	if stats.DistributionScore != 100 {
		t.Fatalf("expected DistributionScore=100 for uniform distribution, got %d", stats.DistributionScore)
	}
	if len(stats.HotPartitions) != 0 {
		t.Fatalf("expected no hot partitions for uniform distribution, got %v", stats.HotPartitions)
	}
}

// TestHotPartitionDetection verifies that a partition holding >10% of items is flagged.
func TestHotPartitionDetection(t *testing.T) {
	items := make([]Item, 0, 100)
	// One partition holds 60 items (60% of total — well above 10% threshold).
	for j := 0; j < 60; j++ {
		items = append(items, Item{PK: "HOT_KEY", SK: "SK#" + string(rune('a'+j%26)), EntityType: "test"})
	}
	// Four other partitions hold 10 items each.
	for i := 1; i <= 4; i++ {
		for j := 0; j < 10; j++ {
			items = append(items, Item{
				PK:         "COLD_KEY_" + string(rune('A'+i)),
				SK:         "SK#" + string(rune('a'+j)),
				EntityType: "test",
			})
		}
	}

	stats := AnalyzeDistribution(items)

	if len(stats.HotPartitions) == 0 {
		t.Fatal("expected HOT_KEY to be detected as a hot partition")
	}

	found := false
	for _, h := range stats.HotPartitions {
		if len(h) >= 7 && h[:7] == "HOT_KEY" {
			found = true
		}
	}
	if !found {
		t.Fatalf("HOT_KEY not found in hot partitions list: %v", stats.HotPartitions)
	}

	if stats.DistributionScore >= 70 {
		t.Fatalf("expected DistributionScore < 70 for heavily skewed distribution, got %d", stats.DistributionScore)
	}
}

// TestTotalItemCount verifies that TotalItems matches the input slice length.
func TestTotalItemCount(t *testing.T) {
	items := []Item{
		{PK: "A", SK: "1", EntityType: "x"},
		{PK: "A", SK: "2", EntityType: "x"},
		{PK: "B", SK: "1", EntityType: "x"},
	}
	stats := AnalyzeDistribution(items)
	if stats.TotalItems != 3 {
		t.Fatalf("expected TotalItems=3, got %d", stats.TotalItems)
	}
}

// TestEmptyInput verifies graceful handling of zero items.
func TestEmptyInput(t *testing.T) {
	stats := AnalyzeDistribution(nil)
	if stats.TotalItems != 0 {
		t.Fatalf("expected TotalItems=0 for empty input, got %d", stats.TotalItems)
	}
	if stats.DistributionScore != 100 {
		t.Fatalf("expected DistributionScore=100 for empty input, got %d", stats.DistributionScore)
	}
}

// TestSinglePartition verifies that a single partition key has a perfect score
// (only one partition, no variance) and is flagged as hot since it holds 100% of items.
func TestSinglePartition(t *testing.T) {
	items := []Item{
		{PK: "ONLY", SK: "1"},
		{PK: "ONLY", SK: "2"},
		{PK: "ONLY", SK: "3"},
	}
	stats := AnalyzeDistribution(items)
	if stats.DistributionScore != 100 {
		t.Fatalf("expected score=100 for single partition (no variance), got %d", stats.DistributionScore)
	}
	if len(stats.HotPartitions) == 0 {
		t.Fatal("expected single partition to be flagged as hot (holds 100% of items)")
	}
}
