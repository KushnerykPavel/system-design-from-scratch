package main

import (
	"encoding/json"
	"flag"
	"os"
)

type AutocompleteConfig struct {
	Name               string `json:"name"`
	MaxSuggestions     int    `json:"max_suggestions"`
	FreshnessSeconds   int    `json:"freshness_seconds"`
	IndexReplicas      int    `json:"index_replicas"`
	PrefixCacheEnabled bool   `json:"prefix_cache_enabled"`
	Personalization    bool   `json:"personalization"`
	PolicyFiltering    bool   `json:"policy_filtering"`
	FallbackSnapshot   bool   `json:"fallback_snapshot"`
}

func ValidateAutocompleteConfig(cfg AutocompleteConfig) []string {
	var issues []string
	if cfg.MaxSuggestions < 3 || cfg.MaxSuggestions > 20 {
		issues = append(issues, "max_suggestions should stay between 3 and 20")
	}
	if cfg.FreshnessSeconds <= 0 || cfg.FreshnessSeconds > 900 {
		issues = append(issues, "freshness_seconds should reflect a realistic trend-refresh target")
	}
	if cfg.IndexReplicas < 2 {
		issues = append(issues, "index_replicas should be at least 2 for serving availability")
	}
	if !cfg.PrefixCacheEnabled {
		issues = append(issues, "prefix_cache_enabled should usually be true for hot prefixes")
	}
	if !cfg.PolicyFiltering {
		issues = append(issues, "policy_filtering should be enabled before serving suggestions")
	}
	if !cfg.FallbackSnapshot {
		issues = append(issues, "fallback_snapshot should be enabled to survive feature-pipeline lag")
	}
	return issues
}

func main() {
	name := flag.String("name", "global-autocomplete", "autocomplete config name")
	flag.Parse()

	cfg := AutocompleteConfig{
		Name:               *name,
		MaxSuggestions:     8,
		FreshnessSeconds:   300,
		IndexReplicas:      3,
		PrefixCacheEnabled: true,
		Personalization:    true,
		PolicyFiltering:    true,
		FallbackSnapshot:   true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateAutocompleteConfig(cfg),
	})
}
