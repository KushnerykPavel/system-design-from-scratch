package main

import (
	"encoding/json"
	"flag"
	"os"
)

type CacheClusterConfig struct {
	Name                string `json:"name"`
	Shards              int    `json:"shards"`
	Replicas            int    `json:"replicas"`
	EvictionPolicy      string `json:"eviction_policy"`
	TargetHitRate       int    `json:"target_hit_rate"`
	RequestCoalescing   bool   `json:"request_coalescing"`
	ConsistentHashing   bool   `json:"consistent_hashing"`
	NegativeCaching     bool   `json:"negative_caching"`
	TenantIsolationMode string `json:"tenant_isolation_mode"`
}

func ValidateCacheCluster(cfg CacheClusterConfig) []string {
	var issues []string
	if cfg.Shards <= 0 {
		issues = append(issues, "shards must be positive")
	}
	if cfg.Replicas < 1 {
		issues = append(issues, "replicas must be at least 1")
	}
	if cfg.EvictionPolicy != "lru" && cfg.EvictionPolicy != "lfu" && cfg.EvictionPolicy != "size_aware" {
		issues = append(issues, "eviction_policy must be lru, lfu, or size_aware")
	}
	if cfg.TargetHitRate < 50 || cfg.TargetHitRate > 99 {
		issues = append(issues, "target_hit_rate should be between 50 and 99 percent")
	}
	if !cfg.RequestCoalescing {
		issues = append(issues, "request_coalescing should be enabled to avoid miss storms")
	}
	if !cfg.ConsistentHashing {
		issues = append(issues, "consistent_hashing should be enabled to reduce key churn on rebalance")
	}
	if cfg.TenantIsolationMode != "shared" && cfg.TenantIsolationMode != "segmented" && cfg.TenantIsolationMode != "dedicated" {
		issues = append(issues, "tenant_isolation_mode must be shared, segmented, or dedicated")
	}
	return issues
}

func main() {
	name := flag.String("name", "shared-cache-cluster", "name of the cache cluster")
	flag.Parse()

	cfg := CacheClusterConfig{
		Name:                *name,
		Shards:              64,
		Replicas:            1,
		EvictionPolicy:      "lfu",
		TargetHitRate:       92,
		RequestCoalescing:   true,
		ConsistentHashing:   true,
		NegativeCaching:     true,
		TenantIsolationMode: "segmented",
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"cluster": cfg,
		"issues":  ValidateCacheCluster(cfg),
	})
}
