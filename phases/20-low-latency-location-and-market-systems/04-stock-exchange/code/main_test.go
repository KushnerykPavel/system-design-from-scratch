package main

import "testing"

func TestValidateMatchingEngineConfigHealthy(t *testing.T) {
	cfg := MatchingEngineConfig{
		Name:                   "healthy",
		SymbolPartitions:       64,
		HotSymbolIsolation:     true,
		JournalReplicas:        3,
		AckAfterDurableJournal: true,
		MarketDataAsyncFanout:  true,
		SnapshotIntervalEvents: 100000,
	}
	if issues := ValidateMatchingEngineConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateMatchingEngineConfig returned issues: %v", issues)
	}
}

func TestValidateMatchingEngineConfigWeak(t *testing.T) {
	cfg := MatchingEngineConfig{
		Name:                   "weak",
		SymbolPartitions:       1,
		JournalReplicas:        1,
		SnapshotIntervalEvents: 10,
	}
	if issues := ValidateMatchingEngineConfig(cfg); len(issues) < 5 {
		t.Fatalf("ValidateMatchingEngineConfig returned too few issues: %v", issues)
	}
}
