package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ReplayRequest struct {
	Name             string `json:"name"`
	HasReason        bool   `json:"has_reason"`
	HasScopedTarget  bool   `json:"has_scoped_target"`
	HasTimeWindow    bool   `json:"has_time_window"`
	HasDryRun        bool   `json:"has_dry_run"`
	HasRateLimit     bool   `json:"has_rate_limit"`
	HasActorIdentity bool   `json:"has_actor_identity"`
	HasRollbackPlan  bool   `json:"has_rollback_plan"`
}

func ValidateReplayRequest(req ReplayRequest) []string {
	var issues []string
	if !req.HasReason {
		issues = append(issues, "has_reason should be true so replay intent is auditable")
	}
	if !req.HasScopedTarget {
		issues = append(issues, "has_scoped_target should be true so replay does not accidentally touch everything")
	}
	if !req.HasTimeWindow {
		issues = append(issues, "has_time_window should be true so blast radius stays bounded")
	}
	if !req.HasDryRun {
		issues = append(issues, "has_dry_run should be true so impact can be estimated before execution")
	}
	if !req.HasRateLimit {
		issues = append(issues, "has_rate_limit should be true so replay cannot starve live traffic")
	}
	if !req.HasActorIdentity {
		issues = append(issues, "has_actor_identity should be true for accountability and incident review")
	}
	if !req.HasRollbackPlan {
		issues = append(issues, "has_rollback_plan should be true when replay could trigger harmful side effects")
	}
	return issues
}

func main() {
	name := flag.String("name", "replay-request", "request name")
	flag.Parse()

	req := ReplayRequest{
		Name:             *name,
		HasReason:        true,
		HasScopedTarget:  true,
		HasTimeWindow:    true,
		HasDryRun:        true,
		HasRateLimit:     true,
		HasActorIdentity: true,
		HasRollbackPlan:  true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"request": req,
		"issues":  ValidateReplayRequest(req),
	})
}
