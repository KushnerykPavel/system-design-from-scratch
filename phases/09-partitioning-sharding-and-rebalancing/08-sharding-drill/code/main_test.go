package main

import "testing"

func TestScoreCountsCoveredTopics(t *testing.T) {
	answer := "Shard key, tenant isolation, directory, rebalancing, observability, and trade-off."
	got := score(answer)
	if got < 6 {
		t.Fatalf("expected at least 6 covered topics, got %d", got)
	}
}

func TestCoverageDetectsMissingTopics(t *testing.T) {
	answer := "We shard by tenant and ignore the rest."
	got := coverage(answer)
	if got["rebalancing"] {
		t.Fatal("expected rebalancing to be missing")
	}
	if !got["shard key"] {
		t.Fatal("expected shard key to be covered")
	}
}
