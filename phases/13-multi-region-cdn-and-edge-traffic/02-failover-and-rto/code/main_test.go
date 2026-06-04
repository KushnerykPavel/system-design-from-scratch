package main

import "testing"

func TestValidateRecoveryPlanAcceptsHealthyPlan(t *testing.T) {
	plan := RecoveryPlan{
		Name:                    "good",
		DetectionSeconds:        30,
		PromotionSeconds:        60,
		TrafficShiftSeconds:     60,
		WarmupSeconds:           30,
		TargetRTOMinutes:        5,
		TargetRPOSeconds:        20,
		ReplicaLagSeconds:       10,
		HasApprovalGate:         true,
		HasDrillEvidence:        true,
		SupportsReadOnlyDegrade: true,
	}
	if issues := ValidateRecoveryPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateRecoveryPlan() returned issues: %v", issues)
	}
}

func TestValidateRecoveryPlanRejectsWeakPlan(t *testing.T) {
	plan := RecoveryPlan{
		Name:                    "bad",
		DetectionSeconds:        120,
		PromotionSeconds:        120,
		TrafficShiftSeconds:     120,
		WarmupSeconds:           120,
		TargetRTOMinutes:        5,
		TargetRPOSeconds:        10,
		ReplicaLagSeconds:       60,
		HasApprovalGate:         false,
		HasDrillEvidence:        false,
		SupportsReadOnlyDegrade: false,
	}
	if issues := ValidateRecoveryPlan(plan); len(issues) < 4 {
		t.Fatalf("ValidateRecoveryPlan() returned too few issues: %v", issues)
	}
}
