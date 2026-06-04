package main

import "testing"

func TestValidateTopologyAcceptsHealthyConfig(t *testing.T) {
	cfg := KVTopology{
		Name:              "healthy",
		Replicas:          3,
		WriteQuorum:       2,
		ReadQuorum:        2,
		FailureDomains:    3,
		ConsistencyMode:   "quorum",
		RepairEnabled:     true,
		HintedHandoff:     true,
		HotKeyMitigation:  true,
		ConditionalWrites: true,
	}
	if issues := ValidateTopology(cfg); len(issues) != 0 {
		t.Fatalf("ValidateTopology returned issues: %v", issues)
	}
}

func TestValidateTopologyRejectsWeakConfig(t *testing.T) {
	cfg := KVTopology{
		Name:            "weak",
		Replicas:        2,
		WriteQuorum:     3,
		ReadQuorum:      0,
		FailureDomains:  1,
		ConsistencyMode: "unknown",
	}
	if issues := ValidateTopology(cfg); len(issues) < 5 {
		t.Fatalf("ValidateTopology returned too few issues: %v", issues)
	}
}
