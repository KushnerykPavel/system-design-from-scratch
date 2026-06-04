package main

import (
	"encoding/json"
	"flag"
	"os"
)

type MatchingEngineConfig struct {
	Name                   string `json:"name"`
	SymbolPartitions       int    `json:"symbol_partitions"`
	HotSymbolIsolation     bool   `json:"hot_symbol_isolation"`
	JournalReplicas        int    `json:"journal_replicas"`
	AckAfterDurableJournal bool   `json:"ack_after_durable_journal"`
	MarketDataAsyncFanout  bool   `json:"market_data_async_fanout"`
	SnapshotIntervalEvents int    `json:"snapshot_interval_events"`
}

func ValidateMatchingEngineConfig(cfg MatchingEngineConfig) []string {
	var issues []string
	if cfg.SymbolPartitions < 4 {
		issues = append(issues, "symbol_partitions should be at least 4")
	}
	if !cfg.HotSymbolIsolation {
		issues = append(issues, "hot_symbol_isolation should usually be enabled")
	}
	if cfg.JournalReplicas < 2 {
		issues = append(issues, "journal_replicas should be at least 2")
	}
	if !cfg.AckAfterDurableJournal {
		issues = append(issues, "ack_after_durable_journal should be true for a credible recovery boundary")
	}
	if !cfg.MarketDataAsyncFanout {
		issues = append(issues, "market_data_async_fanout should be enabled to protect the matching path")
	}
	if cfg.SnapshotIntervalEvents < 1000 || cfg.SnapshotIntervalEvents > 10000000 {
		issues = append(issues, "snapshot_interval_events should stay between 1000 and 10000000")
	}
	return issues
}

func main() {
	name := flag.String("name", "matching-engine", "config name")
	flag.Parse()

	cfg := MatchingEngineConfig{
		Name:                   *name,
		SymbolPartitions:       64,
		HotSymbolIsolation:     true,
		JournalReplicas:        3,
		AckAfterDurableJournal: true,
		MarketDataAsyncFanout:  true,
		SnapshotIntervalEvents: 100000,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateMatchingEngineConfig(cfg),
	})
}
