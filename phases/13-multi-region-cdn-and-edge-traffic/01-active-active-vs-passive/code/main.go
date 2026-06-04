package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type TopologyProfile struct {
	Name                    string `json:"name"`
	Mode                    string `json:"mode"`
	Regions                 int    `json:"regions"`
	ConcurrentWrites        bool   `json:"concurrent_writes"`
	ReservedFailoverPercent int    `json:"reserved_failover_percent"`
	AutomatedFailover       bool   `json:"automated_failover"`
	HasReadinessChecks      bool   `json:"has_readiness_checks"`
	HasFailbackPlan         bool   `json:"has_failback_plan"`
}

func ValidateTopology(profile TopologyProfile) []string {
	var issues []string
	if profile.Mode != "active_active" && profile.Mode != "active_passive" {
		issues = append(issues, "mode must be active_active or active_passive")
	}
	if profile.Regions < 2 {
		issues = append(issues, "at least two regions are required")
	}
	if profile.Mode == "active_active" && profile.ConcurrentWrites && !profile.HasFailbackPlan {
		issues = append(issues, "multi-writer active-active requires an explicit failback or reconciliation plan")
	}
	if profile.ReservedFailoverPercent < 20 {
		issues = append(issues, "reserved failover capacity should be at least 20 percent")
	}
	if !profile.AutomatedFailover {
		issues = append(issues, "regional recovery should not depend entirely on manual failover")
	}
	if !profile.HasReadinessChecks {
		issues = append(issues, "passive or peer regions need explicit readiness checks")
	}
	return issues
}

func main() {
	name := flag.String("name", "global-service", "topology name")
	flag.Parse()

	profile := TopologyProfile{
		Name:                    *name,
		Mode:                    "active_passive",
		Regions:                 2,
		ConcurrentWrites:        false,
		ReservedFailoverPercent: 35,
		AutomatedFailover:       true,
		HasReadinessChecks:      true,
		HasFailbackPlan:         true,
	}

	payload := map[string]any{
		"profile": profile,
		"issues":  ValidateTopology(profile),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
