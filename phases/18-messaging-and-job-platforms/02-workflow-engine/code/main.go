package main

import (
	"encoding/json"
	"flag"
	"os"
)

type WorkflowPlan struct {
	Name                    string `json:"name"`
	HasDurableState         bool   `json:"has_durable_state"`
	HasTimerQueue           bool   `json:"has_timer_queue"`
	HasIdempotentActivities bool   `json:"has_idempotent_activities"`
	SupportsCompensation    bool   `json:"supports_compensation"`
	HasHeartbeats           bool   `json:"has_heartbeats"`
	HasVersioning           bool   `json:"has_versioning"`
	HasIsolationControls    bool   `json:"has_isolation_controls"`
}

func ValidateWorkflowPlan(plan WorkflowPlan) []string {
	var issues []string
	if !plan.HasDurableState {
		issues = append(issues, "has_durable_state should be true so workflow progress survives worker failure")
	}
	if !plan.HasTimerQueue {
		issues = append(issues, "has_timer_queue should be true for sleeps, deadlines, and delayed retries")
	}
	if !plan.HasIdempotentActivities {
		issues = append(issues, "has_idempotent_activities should be true because activity retries are normal")
	}
	if !plan.SupportsCompensation {
		issues = append(issues, "supports_compensation should be true when long-running workflows coordinate side effects")
	}
	if !plan.HasHeartbeats {
		issues = append(issues, "has_heartbeats should be true so stuck work can be detected")
	}
	if !plan.HasVersioning {
		issues = append(issues, "has_versioning should be true so in-flight workflows survive code evolution")
	}
	if !plan.HasIsolationControls {
		issues = append(issues, "has_isolation_controls should be true so one workflow type does not starve others")
	}
	return issues
}

func main() {
	name := flag.String("name", "workflow-engine", "plan name")
	flag.Parse()

	plan := WorkflowPlan{
		Name:                    *name,
		HasDurableState:         true,
		HasTimerQueue:           true,
		HasIdempotentActivities: true,
		SupportsCompensation:    true,
		HasHeartbeats:           true,
		HasVersioning:           true,
		HasIsolationControls:    true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateWorkflowPlan(plan),
	})
}
