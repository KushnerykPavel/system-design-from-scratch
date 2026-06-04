package main

import (
	"encoding/json"
	"flag"
	"os"
)

type CollaborationPlan struct {
	Name                  string `json:"name"`
	ConvergenceMode       string `json:"convergence_mode"`
	HasSessionOwner       bool   `json:"has_session_owner"`
	HasSnapshots          bool   `json:"has_snapshots"`
	SnapshotIntervalOps   int    `json:"snapshot_interval_ops"`
	HasPresenceTTL        bool   `json:"has_presence_ttl"`
	SupportsReplay        bool   `json:"supports_replay"`
	HasDeterministicCheck bool   `json:"has_deterministic_check"`
}

func ValidateCollaborationPlan(plan CollaborationPlan) []string {
	var issues []string
	if plan.ConvergenceMode != "ot" && plan.ConvergenceMode != "crdt" && plan.ConvergenceMode != "sequenced" {
		issues = append(issues, "convergence_mode must be ot, crdt, or sequenced")
	}
	if !plan.HasSessionOwner && plan.ConvergenceMode == "sequenced" {
		issues = append(issues, "has_session_owner should be true for sequenced coordination")
	}
	if !plan.HasSnapshots {
		issues = append(issues, "has_snapshots should be true to cap replay cost")
	}
	if plan.SnapshotIntervalOps <= 0 {
		issues = append(issues, "snapshot_interval_ops must be positive")
	}
	if !plan.HasPresenceTTL {
		issues = append(issues, "has_presence_ttl should be true to clear ghost collaborators")
	}
	if !plan.SupportsReplay {
		issues = append(issues, "supports_replay should be true for reconnect recovery")
	}
	if !plan.HasDeterministicCheck {
		issues = append(issues, "has_deterministic_check should be true to detect corruption in merge logic")
	}
	return issues
}

func main() {
	name := flag.String("name", "doc-collab", "plan name")
	flag.Parse()

	plan := CollaborationPlan{
		Name:                  *name,
		ConvergenceMode:       "sequenced",
		HasSessionOwner:       true,
		HasSnapshots:          true,
		SnapshotIntervalOps:   500,
		HasPresenceTTL:        true,
		SupportsReplay:        true,
		HasDeterministicCheck: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateCollaborationPlan(plan),
	})
}
