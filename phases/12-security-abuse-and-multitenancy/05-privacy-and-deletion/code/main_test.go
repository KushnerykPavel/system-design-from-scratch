package main

import "testing"

func TestValidateDeletionPlan(t *testing.T) {
	if issues := ValidateDeletionPlan(DeletionPlan{
		ImmediateHide:        true,
		AsyncFanout:          true,
		Tombstones:           true,
		BackupRecoveryReplay: true,
	}); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	if issues := ValidateDeletionPlan(DeletionPlan{ClaimsInstantBackupPurge: true}); len(issues) < 4 {
		t.Fatalf("expected several issues, got %v", issues)
	}
}
