package main

import (
	"encoding/json"
	"flag"
	"os"
)

type RecoveryPlan struct {
	Name                  string `json:"name"`
	ExplicitStates        bool   `json:"explicit_states"`
	IdempotentTransitions bool   `json:"idempotent_transitions"`
	TimeoutStates         bool   `json:"timeout_states"`
	CompensationPath      bool   `json:"compensation_path"`
	EventHistory          bool   `json:"event_history"`
	OperatorRecovery      bool   `json:"operator_recovery"`
	CallbackDedupe        bool   `json:"callback_dedupe"`
}

func ValidateRecoveryPlan(plan RecoveryPlan) []string {
	var issues []string
	if !plan.ExplicitStates {
		issues = append(issues, "explicit_states should be true so ambiguous workflow phases are modeled directly")
	}
	if !plan.IdempotentTransitions {
		issues = append(issues, "idempotent_transitions should be true so repeated callbacks stay safe")
	}
	if !plan.TimeoutStates {
		issues = append(issues, "timeout_states should be true so unknown external outcomes are not hidden")
	}
	if !plan.CompensationPath {
		issues = append(issues, "compensation_path should be true so partial failures have a recovery plan")
	}
	if !plan.EventHistory {
		issues = append(issues, "event_history should be true so order debugging and replay are possible")
	}
	if !plan.OperatorRecovery {
		issues = append(issues, "operator_recovery should be true so stalled orders can be resumed safely")
	}
	if !plan.CallbackDedupe {
		issues = append(issues, "callback_dedupe should be true so out-of-order or repeated external events do not corrupt state")
	}
	return issues
}

func main() {
	name := flag.String("name", "order-state-machine", "plan name")
	flag.Parse()

	plan := RecoveryPlan{
		Name:                  *name,
		ExplicitStates:        true,
		IdempotentTransitions: true,
		TimeoutStates:         true,
		CompensationPath:      true,
		EventHistory:          true,
		OperatorRecovery:      true,
		CallbackDedupe:        true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateRecoveryPlan(plan),
	})
}
