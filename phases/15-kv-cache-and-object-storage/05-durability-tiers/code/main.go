package main

import (
	"encoding/json"
	"flag"
	"os"
)

type DurabilityTier struct {
	Name             string `json:"name"`
	ReplicaCopies    int    `json:"replica_copies"`
	ErasureDataParts int    `json:"erasure_data_parts"`
	ErasureParity    int    `json:"erasure_parity"`
	GeoRedundant     bool   `json:"geo_redundant"`
	MaxRepairHours   int    `json:"max_repair_hours"`
	RestoreTested    bool   `json:"restore_tested"`
}

func ValidateDurabilityTier(cfg DurabilityTier) []string {
	var issues []string
	if cfg.ReplicaCopies == 0 && (cfg.ErasureDataParts == 0 || cfg.ErasureParity == 0) {
		issues = append(issues, "tier must define either replica copies or an erasure-coding scheme")
	}
	if cfg.ReplicaCopies > 0 && cfg.ReplicaCopies < 3 && !cfg.GeoRedundant {
		issues = append(issues, "regional replicated tiers usually need at least 3 copies")
	}
	if cfg.MaxRepairHours <= 0 {
		issues = append(issues, "max_repair_hours must be positive")
	}
	if !cfg.RestoreTested {
		issues = append(issues, "restore_tested should be true for any tier with durability claims")
	}
	if cfg.ErasureDataParts > 0 && cfg.ErasureParity <= 0 {
		issues = append(issues, "erasure-coded tiers need positive parity")
	}
	return issues
}

func main() {
	name := flag.String("name", "standard", "name of the durability tier")
	flag.Parse()

	cfg := DurabilityTier{
		Name:           *name,
		ReplicaCopies:  3,
		GeoRedundant:   false,
		MaxRepairHours: 6,
		RestoreTested:  true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"tier":   cfg,
		"issues": ValidateDurabilityTier(cfg),
	})
}
