package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type RingPlan struct {
	Name                    string `json:"name"`
	Nodes                   int    `json:"nodes"`
	VirtualNodesPerNode     int    `json:"virtual_nodes_per_node"`
	ExpectedRemapPercent    int    `json:"expected_remap_percent"`
	WeightedPlacement       bool   `json:"weighted_placement"`
	HealthAwareRouting      bool   `json:"health_aware_routing"`
	PrewarmOnRebalance      bool   `json:"prewarm_on_rebalance"`
	HotKeyIsolationStrategy bool   `json:"hot_key_isolation_strategy"`
}

func ValidateRingPlan(plan RingPlan) []string {
	var issues []string
	if plan.Nodes < 2 {
		issues = append(issues, "at least two nodes are required")
	}
	if plan.VirtualNodesPerNode < 16 {
		issues = append(issues, "virtual_nodes_per_node should usually be at least 16 for smoother balance")
	}
	if plan.ExpectedRemapPercent > 35 {
		issues = append(issues, "expected_remap_percent looks high for a bounded rebalance plan")
	}
	if !plan.HealthAwareRouting {
		issues = append(issues, "ring placement needs health-aware routing to avoid dead owners")
	}
	if !plan.PrewarmOnRebalance {
		issues = append(issues, "rebalance should plan for prewarming or staged migration")
	}
	if !plan.HotKeyIsolationStrategy {
		issues = append(issues, "consistent hashing alone does not solve pathological hot keys")
	}
	return issues
}

func main() {
	name := flag.String("name", "edge-cache-ring", "ring plan name")
	flag.Parse()

	plan := RingPlan{
		Name:                    *name,
		Nodes:                   100,
		VirtualNodesPerNode:     64,
		ExpectedRemapPercent:    12,
		WeightedPlacement:       true,
		HealthAwareRouting:      true,
		PrewarmOnRebalance:      true,
		HotKeyIsolationStrategy: true,
	}

	payload := map[string]any{
		"plan":   plan,
		"issues": ValidateRingPlan(plan),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
