package main

import (
	"encoding/json"
	"flag"
	"os"
)

type InfraPlatformScorecard struct {
	Name                         string `json:"name"`
	ClarifiesPlatformContract    bool   `json:"clarifies_platform_contract"`
	SeparatesControlAndDataPlane bool   `json:"separates_control_and_data_plane"`
	SizesHotAndSlowPaths         bool   `json:"sizes_hot_and_slow_paths"`
	DefinesFailureDomain         bool   `json:"defines_failure_domain"`
	HasRolloutOrMigrationPlan    bool   `json:"has_rollout_or_migration_plan"`
	CoversTenantGuardrails       bool   `json:"covers_tenant_guardrails"`
	CoversObservability          bool   `json:"covers_observability"`
	HandlesConstraintChange      bool   `json:"handles_constraint_change"`
}

func ValidateInfraPlatform(card InfraPlatformScorecard) []string {
	var issues []string
	if !card.ClarifiesPlatformContract {
		issues = append(issues, "clarifies_platform_contract should be true so the answer defines who the platform serves")
	}
	if !card.SeparatesControlAndDataPlane {
		issues = append(issues, "separates_control_and_data_plane should be true so safety and latency are discussed on the right paths")
	}
	if !card.SizesHotAndSlowPaths {
		issues = append(issues, "sizes_hot_and_slow_paths should be true so both data-plane and control-plane choices are grounded")
	}
	if !card.DefinesFailureDomain {
		issues = append(issues, "defines_failure_domain should be true so blast radius is explicit")
	}
	if !card.HasRolloutOrMigrationPlan {
		issues = append(issues, "has_rollout_or_migration_plan should be true so state propagation is not magical")
	}
	if !card.CoversTenantGuardrails {
		issues = append(issues, "covers_tenant_guardrails should be true so noisy-neighbor and misconfiguration risks are contained")
	}
	if !card.CoversObservability {
		issues = append(issues, "covers_observability should be true so propagation and serving behavior are measurable")
	}
	if !card.HandlesConstraintChange {
		issues = append(issues, "handles_constraint_change should be true so new rollout or latency demands change the design")
	}
	return issues
}

func main() {
	name := flag.String("name", "infra-platform-mock", "mock name")
	flag.Parse()

	card := InfraPlatformScorecard{
		Name:                         *name,
		ClarifiesPlatformContract:    true,
		SeparatesControlAndDataPlane: true,
		SizesHotAndSlowPaths:         true,
		DefinesFailureDomain:         true,
		HasRolloutOrMigrationPlan:    true,
		CoversTenantGuardrails:       true,
		CoversObservability:          true,
		HandlesConstraintChange:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateInfraPlatform(card),
	})
}
