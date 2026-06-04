package main

import "testing"

func TestRankPatterns(t *testing.T) {
	patterns := []AccessPattern{
		{Name: "export", ReadQPS: 50, WriteQPS: 0, LatencyCritical: false, RequiresStrongTxn: false, SupportsRevenue: false},
		{Name: "checkout", ReadQPS: 200, WriteQPS: 80, LatencyCritical: true, RequiresStrongTxn: true, SupportsRevenue: true},
		{Name: "search", ReadQPS: 800, WriteQPS: 10, LatencyCritical: false, RequiresStrongTxn: false, SupportsRevenue: false},
	}

	ranked := RankPatterns(patterns)
	if len(ranked) != 3 {
		t.Fatalf("expected 3 ranked patterns, got %d", len(ranked))
	}
	if ranked[0].Name != "checkout" {
		t.Fatalf("expected checkout to rank first, got %s", ranked[0].Name)
	}
	if ranked[2].Name != "export" {
		t.Fatalf("expected export to rank last, got %s", ranked[2].Name)
	}
}
