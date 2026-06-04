package main

import "testing"

func TestAssessIsolationFlagsSharedBlastRadius(t *testing.T) {
	got := AssessIsolation(IsolationPlan{
		DedicatedCriticalPool: false,
		TenantQuota:           false,
		CellCount:             1,
		SharedCriticalDep:     true,
		AllowsBorrowing:       true,
		LargestTenantShare:    0.18,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessIsolationApprovesBoundedPlan(t *testing.T) {
	got := AssessIsolation(IsolationPlan{
		DedicatedCriticalPool: true,
		TenantQuota:           true,
		CellCount:             8,
		SharedCriticalDep:     false,
		AllowsBorrowing:       true,
		LargestTenantShare:    0.15,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
