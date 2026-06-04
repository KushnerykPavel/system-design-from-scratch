package main

import (
	"encoding/json"
	"os"
)

type DeletionPlan struct {
	ImmediateHide            bool `json:"immediate_hide"`
	AsyncFanout              bool `json:"async_fanout"`
	Tombstones               bool `json:"tombstones"`
	BackupRecoveryReplay     bool `json:"backup_recovery_replay"`
	ClaimsInstantBackupPurge bool `json:"claims_instant_backup_purge"`
}

func ValidateDeletionPlan(plan DeletionPlan) []string {
	var issues []string
	if !plan.ImmediateHide {
		issues = append(issues, "user-facing reads should be suppressed quickly after delete request")
	}
	if !plan.AsyncFanout {
		issues = append(issues, "derived systems need asynchronous deletion fanout")
	}
	if !plan.Tombstones {
		issues = append(issues, "tombstones or deletion ledger help prevent resurrection")
	}
	if !plan.BackupRecoveryReplay {
		issues = append(issues, "restores should replay deletion state")
	}
	if plan.ClaimsInstantBackupPurge {
		issues = append(issues, "do not promise instant purge from immutable backups")
	}
	return issues
}

func main() {
	plan := DeletionPlan{
		ImmediateHide:        true,
		AsyncFanout:          true,
		Tombstones:           true,
		BackupRecoveryReplay: true,
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateDeletionPlan(plan),
	})
}
