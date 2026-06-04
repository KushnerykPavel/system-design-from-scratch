package main

import "testing"

func TestValidateCacheClusterAcceptsReasonableConfig(t *testing.T) {
	cfg := CacheClusterConfig{
		Name:                "reasonable",
		Shards:              64,
		Replicas:            1,
		EvictionPolicy:      "lfu",
		TargetHitRate:       90,
		RequestCoalescing:   true,
		ConsistentHashing:   true,
		NegativeCaching:     true,
		TenantIsolationMode: "segmented",
	}
	if issues := ValidateCacheCluster(cfg); len(issues) != 0 {
		t.Fatalf("ValidateCacheCluster returned issues: %v", issues)
	}
}

func TestValidateCacheClusterRejectsRiskySettings(t *testing.T) {
	cfg := CacheClusterConfig{
		Name:                "risky",
		Shards:              0,
		Replicas:            0,
		EvictionPolicy:      "fifo",
		TargetHitRate:       20,
		RequestCoalescing:   false,
		ConsistentHashing:   false,
		TenantIsolationMode: "none",
	}
	if issues := ValidateCacheCluster(cfg); len(issues) < 5 {
		t.Fatalf("ValidateCacheCluster returned too few issues: %v", issues)
	}
}
