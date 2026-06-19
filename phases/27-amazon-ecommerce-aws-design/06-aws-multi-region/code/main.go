package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// DRTier represents a disaster recovery tier.
type DRTier string

const (
	// BackupRestore: RTO > 60 min, RPO > 60 min. S3 backup + restore on demand.
	BackupRestore DRTier = "BACKUP_RESTORE"
	// PilotLight: RTO 15–60 min, RPO 5–60 min. Minimal infrastructure pre-deployed.
	PilotLight DRTier = "PILOT_LIGHT"
	// WarmStandby: RTO 1–15 min, RPO 1–5 min. Pre-scaled standby, not serving traffic.
	WarmStandby DRTier = "WARM_STANDBY"
	// MultiActive: RTO < 1 min, RPO < 1 min. All regions serve live traffic.
	MultiActive DRTier = "MULTI_ACTIVE"
)

// DRRequirements specifies the RTO and RPO constraints for a workload.
type DRRequirements struct {
	// RTOMinutes is the maximum tolerable downtime in minutes.
	RTOMinutes int
	// RPOMinutes is the maximum tolerable data loss in minutes.
	RPOMinutes int
}

// SelectDRTier returns the minimum DR tier that satisfies the given RTO and RPO requirements.
//
// Tier selection thresholds:
//   - MULTI_ACTIVE:    RTO < 1 min  AND RPO < 1 min
//   - WARM_STANDBY:   RTO ≤ 15 min AND RPO ≤ 5 min
//   - PILOT_LIGHT:    RTO ≤ 60 min AND RPO ≤ 60 min
//   - BACKUP_RESTORE: RTO > 60 min  OR  RPO > 60 min
func SelectDRTier(req DRRequirements) DRTier {
	switch {
	case req.RTOMinutes < 1 && req.RPOMinutes < 1:
		return MultiActive
	case req.RTOMinutes <= 15 && req.RPOMinutes <= 5:
		return WarmStandby
	case req.RTOMinutes <= 60 && req.RPOMinutes <= 60:
		return PilotLight
	default:
		return BackupRestore
	}
}

// DRTierCost returns a relative cost indicator for the given DR tier.
// Cost reflects the additional infrastructure required vs a single-region baseline.
func DRTierCost(tier DRTier) string {
	switch tier {
	case MultiActive:
		return "very-high"
	case WarmStandby:
		return "high"
	case PilotLight:
		return "medium"
	default:
		return "low"
	}
}

// DRResult holds the selected tier, cost, and a human-readable description.
type DRResult struct {
	Tier        DRTier `json:"tier"`
	Cost        string `json:"cost"`
	Description string `json:"description"`
}

func describeTier(tier DRTier, req DRRequirements) string {
	switch tier {
	case MultiActive:
		return fmt.Sprintf("RTO=%dmin RPO=%dmin → MULTI_ACTIVE: all regions serve live traffic; DynamoDB Global Tables multi-master; Global Accelerator anycast routing; near-zero RTO/RPO", req.RTOMinutes, req.RPOMinutes)
	case WarmStandby:
		return fmt.Sprintf("RTO=%dmin RPO=%dmin → WARM_STANDBY: standby region pre-scaled, not serving traffic; Route53 failover routing; Aurora Global Database with automated promotion runbook", req.RTOMinutes, req.RPOMinutes)
	case PilotLight:
		return fmt.Sprintf("RTO=%dmin RPO=%dmin → PILOT_LIGHT: minimal infrastructure (DB replica, cold ECS tasks) in DR region; scale-up on failover; S3 CRR for object storage", req.RTOMinutes, req.RPOMinutes)
	default:
		return fmt.Sprintf("RTO=%dmin RPO=%dmin → BACKUP_RESTORE: nightly S3 backups; restore on demand; no standby infrastructure; acceptable for non-critical workloads", req.RTOMinutes, req.RPOMinutes)
	}
}

func main() {
	testCases := []DRRequirements{
		{RTOMinutes: 0, RPOMinutes: 0},   // MULTI_ACTIVE: both < 1
		{RTOMinutes: 5, RPOMinutes: 1},   // WARM_STANDBY: RTO≤15, RPO≤5
		{RTOMinutes: 15, RPOMinutes: 5},  // WARM_STANDBY: exactly at boundary
		{RTOMinutes: 30, RPOMinutes: 15}, // PILOT_LIGHT: RTO≤60, RPO≤60
		{RTOMinutes: 60, RPOMinutes: 60}, // PILOT_LIGHT: exactly at boundary
		{RTOMinutes: 120, RPOMinutes: 60}, // BACKUP_RESTORE: RTO > 60
		{RTOMinutes: 240, RPOMinutes: 120}, // BACKUP_RESTORE: both > 60
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for _, req := range testCases {
		tier := SelectDRTier(req)
		result := DRResult{
			Tier:        tier,
			Cost:        DRTierCost(tier),
			Description: describeTier(tier, req),
		}
		_ = enc.Encode(result)
	}
}
