package main

import "testing"

func TestValidateTopologyAcceptsReasonableProfile(t *testing.T) {
	profile := TopologyProfile{
		Name:                    "good",
		Mode:                    "active_passive",
		Regions:                 2,
		ConcurrentWrites:        false,
		ReservedFailoverPercent: 30,
		AutomatedFailover:       true,
		HasReadinessChecks:      true,
		HasFailbackPlan:         true,
	}
	if issues := ValidateTopology(profile); len(issues) != 0 {
		t.Fatalf("ValidateTopology() returned issues: %v", issues)
	}
}

func TestValidateTopologyRejectsWeakProfile(t *testing.T) {
	profile := TopologyProfile{
		Name:                    "bad",
		Mode:                    "active_active",
		Regions:                 1,
		ConcurrentWrites:        true,
		ReservedFailoverPercent: 10,
		AutomatedFailover:       false,
		HasReadinessChecks:      false,
		HasFailbackPlan:         false,
	}
	if issues := ValidateTopology(profile); len(issues) < 4 {
		t.Fatalf("ValidateTopology() returned too few issues: %v", issues)
	}
}
