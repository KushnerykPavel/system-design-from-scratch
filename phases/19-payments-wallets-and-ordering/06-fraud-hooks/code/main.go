package main

import (
	"encoding/json"
	"flag"
	"os"
)

type FraudHookPlan struct {
	Name                 string `json:"name"`
	InlineBoundedChecks  bool   `json:"inline_bounded_checks"`
	AsyncScoring         bool   `json:"async_scoring"`
	ManualReviewLane     bool   `json:"manual_review_lane"`
	ExplicitFallback     bool   `json:"explicit_fallback"`
	ModelVersionLogging  bool   `json:"model_version_logging"`
	PolicyRollback       bool   `json:"policy_rollback"`
	OutcomeAuditability  bool   `json:"outcome_auditability"`
}

func ValidateFraudHookPlan(plan FraudHookPlan) []string {
	var issues []string
	if !plan.InlineBoundedChecks {
		issues = append(issues, "inline_bounded_checks should be true so the checkout path has a real latency budget")
	}
	if !plan.AsyncScoring {
		issues = append(issues, "async_scoring should be true so richer risk analysis does not block the core path")
	}
	if !plan.ManualReviewLane {
		issues = append(issues, "manual_review_lane should be true so uncertain outcomes do not collapse into binary allow or deny")
	}
	if !plan.ExplicitFallback {
		issues = append(issues, "explicit_fallback should be true so timeout behavior is deliberate and auditable")
	}
	if !plan.ModelVersionLogging {
		issues = append(issues, "model_version_logging should be true so false positives can be traced to deployed logic")
	}
	if !plan.PolicyRollback {
		issues = append(issues, "policy_rollback should be true so a bad rule or model can be reversed quickly")
	}
	if !plan.OutcomeAuditability {
		issues = append(issues, "outcome_auditability should be true so later disputes can reconstruct why a decision happened")
	}
	return issues
}

func main() {
	name := flag.String("name", "fraud-hooks", "plan name")
	flag.Parse()

	plan := FraudHookPlan{
		Name:                *name,
		InlineBoundedChecks: true,
		AsyncScoring:        true,
		ManualReviewLane:    true,
		ExplicitFallback:    true,
		ModelVersionLogging: true,
		PolicyRollback:      true,
		OutcomeAuditability: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateFraudHookPlan(plan),
	})
}
