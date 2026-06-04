package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type RecoveryPlan struct {
	Name                    string `json:"name"`
	DetectionSeconds        int    `json:"detection_seconds"`
	PromotionSeconds        int    `json:"promotion_seconds"`
	TrafficShiftSeconds     int    `json:"traffic_shift_seconds"`
	WarmupSeconds           int    `json:"warmup_seconds"`
	TargetRTOMinutes        int    `json:"target_rto_minutes"`
	TargetRPOSeconds        int    `json:"target_rpo_seconds"`
	ReplicaLagSeconds       int    `json:"replica_lag_seconds"`
	HasApprovalGate         bool   `json:"has_approval_gate"`
	HasDrillEvidence        bool   `json:"has_drill_evidence"`
	SupportsReadOnlyDegrade bool   `json:"supports_read_only_degrade"`
}

func TotalRecoverySeconds(plan RecoveryPlan) int {
	return plan.DetectionSeconds + plan.PromotionSeconds + plan.TrafficShiftSeconds + plan.WarmupSeconds
}

func ValidateRecoveryPlan(plan RecoveryPlan) []string {
	var issues []string
	if plan.TargetRTOMinutes <= 0 {
		issues = append(issues, "target RTO must be positive")
	}
	if TotalRecoverySeconds(plan) > plan.TargetRTOMinutes*60 {
		issues = append(issues, "recovery steps exceed target RTO")
	}
	if plan.ReplicaLagSeconds > plan.TargetRPOSeconds {
		issues = append(issues, "replica lag exceeds target RPO")
	}
	if !plan.HasApprovalGate {
		issues = append(issues, "failover plan should define an approval or control gate")
	}
	if !plan.HasDrillEvidence {
		issues = append(issues, "failover claims should be backed by drill evidence")
	}
	if !plan.SupportsReadOnlyDegrade {
		issues = append(issues, "consider a read-only degraded mode to reduce outage severity")
	}
	return issues
}

func main() {
	name := flag.String("name", "regional-dr-plan", "recovery plan name")
	flag.Parse()

	plan := RecoveryPlan{
		Name:                    *name,
		DetectionSeconds:        45,
		PromotionSeconds:        60,
		TrafficShiftSeconds:     75,
		WarmupSeconds:           60,
		TargetRTOMinutes:        5,
		TargetRPOSeconds:        30,
		ReplicaLagSeconds:       10,
		HasApprovalGate:         true,
		HasDrillEvidence:        true,
		SupportsReadOnlyDegrade: true,
	}

	payload := map[string]any{
		"plan":                plan,
		"total_recovery_secs": TotalRecoverySeconds(plan),
		"validation_findings": ValidateRecoveryPlan(plan),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
