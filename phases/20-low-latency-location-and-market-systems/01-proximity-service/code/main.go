package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ProximityConfig struct {
	Name                 string `json:"name"`
	DefaultRadiusM       int    `json:"default_radius_m"`
	MaxRadiusM           int    `json:"max_radius_m"`
	MaxCandidates        int    `json:"max_candidates"`
	HotCellReplicas      int    `json:"hot_cell_replicas"`
	LocationTTLSeconds   int    `json:"location_ttl_seconds"`
	AvailabilityCacheTTL int    `json:"availability_cache_ttl_seconds"`
	ExactRerankEnabled   bool   `json:"exact_rerank_enabled"`
}

func ValidateProximityConfig(cfg ProximityConfig) []string {
	var issues []string
	if cfg.DefaultRadiusM <= 0 || cfg.DefaultRadiusM > cfg.MaxRadiusM {
		issues = append(issues, "default_radius_m must be positive and not exceed max_radius_m")
	}
	if cfg.MaxRadiusM < 1000 || cfg.MaxRadiusM > 50000 {
		issues = append(issues, "max_radius_m should stay between 1000 and 50000")
	}
	if cfg.MaxCandidates < 20 || cfg.MaxCandidates > 5000 {
		issues = append(issues, "max_candidates should stay between 20 and 5000")
	}
	if cfg.HotCellReplicas < 2 {
		issues = append(issues, "hot_cell_replicas should be at least 2 for hotspot resilience")
	}
	if cfg.LocationTTLSeconds <= 0 || cfg.LocationTTLSeconds > 120 {
		issues = append(issues, "location_ttl_seconds should reflect a realistic freshness window")
	}
	if cfg.AvailabilityCacheTTL <= 0 || cfg.AvailabilityCacheTTL > 60 {
		issues = append(issues, "availability_cache_ttl_seconds should stay short for moving supply")
	}
	if !cfg.ExactRerankEnabled {
		issues = append(issues, "exact_rerank_enabled should usually be true to avoid cell-boundary errors")
	}
	return issues
}

func main() {
	name := flag.String("name", "nearby-search", "config name")
	flag.Parse()

	cfg := ProximityConfig{
		Name:                 *name,
		DefaultRadiusM:       1000,
		MaxRadiusM:           10000,
		MaxCandidates:        400,
		HotCellReplicas:      3,
		LocationTTLSeconds:   15,
		AvailabilityCacheTTL: 5,
		ExactRerankEnabled:   true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateProximityConfig(cfg),
	})
}
