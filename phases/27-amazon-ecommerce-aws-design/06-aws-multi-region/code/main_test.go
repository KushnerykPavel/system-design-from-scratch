package main

import "testing"

// TestMultiActiveSelection verifies RTO=0, RPO=0 selects MULTI_ACTIVE.
func TestMultiActiveSelection(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 0, RPOMinutes: 0})
	if tier != MultiActive {
		t.Fatalf("expected MULTI_ACTIVE for RTO=0 RPO=0, got %s", tier)
	}
}

// TestWarmStandbySelection verifies RTO=5, RPO=1 selects WARM_STANDBY.
func TestWarmStandbySelection(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 5, RPOMinutes: 1})
	if tier != WarmStandby {
		t.Fatalf("expected WARM_STANDBY for RTO=5 RPO=1, got %s", tier)
	}
}

// TestWarmStandbyBoundary verifies the exact boundary values RTO=15, RPO=5 select WARM_STANDBY.
func TestWarmStandbyBoundary(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 15, RPOMinutes: 5})
	if tier != WarmStandby {
		t.Fatalf("expected WARM_STANDBY for RTO=15 RPO=5, got %s", tier)
	}
}

// TestPilotLightSelection verifies RTO=30, RPO=15 selects PILOT_LIGHT.
func TestPilotLightSelection(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 30, RPOMinutes: 15})
	if tier != PilotLight {
		t.Fatalf("expected PILOT_LIGHT for RTO=30 RPO=15, got %s", tier)
	}
}

// TestPilotLightBoundary verifies the exact boundary values RTO=60, RPO=60 select PILOT_LIGHT.
func TestPilotLightBoundary(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 60, RPOMinutes: 60})
	if tier != PilotLight {
		t.Fatalf("expected PILOT_LIGHT for RTO=60 RPO=60, got %s", tier)
	}
}

// TestBackupRestoreHighRTO verifies RTO > 60 selects BACKUP_RESTORE.
func TestBackupRestoreHighRTO(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 120, RPOMinutes: 30})
	if tier != BackupRestore {
		t.Fatalf("expected BACKUP_RESTORE for RTO=120, got %s", tier)
	}
}

// TestBackupRestoreHighRPO verifies RPO > 60 selects BACKUP_RESTORE even if RTO is low.
func TestBackupRestoreHighRPO(t *testing.T) {
	tier := SelectDRTier(DRRequirements{RTOMinutes: 30, RPOMinutes: 120})
	if tier != BackupRestore {
		t.Fatalf("expected BACKUP_RESTORE for RPO=120, got %s", tier)
	}
}

// TestDRTierCost verifies cost levels are assigned correctly.
func TestDRTierCost(t *testing.T) {
	cases := []struct {
		tier     DRTier
		expected string
	}{
		{MultiActive, "very-high"},
		{WarmStandby, "high"},
		{PilotLight, "medium"},
		{BackupRestore, "low"},
	}
	for _, tc := range cases {
		got := DRTierCost(tc.tier)
		if got != tc.expected {
			t.Errorf("DRTierCost(%s): expected %s, got %s", tc.tier, tc.expected, got)
		}
	}
}
