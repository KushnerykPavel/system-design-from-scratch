package main

import (
	"encoding/json"
	"flag"
	"os"
)

type LedgerPlan struct {
	Name                    string `json:"name"`
	AppendOnly              bool   `json:"append_only"`
	DoubleEntry             bool   `json:"double_entry"`
	IdempotentWrites        bool   `json:"idempotent_writes"`
	ImmutableAuditTrail     bool   `json:"immutable_audit_trail"`
	ReconciliationWorkflow  bool   `json:"reconciliation_workflow"`
	ProjectionRebuild       bool   `json:"projection_rebuild"`
	CorrectionByReversal    bool   `json:"correction_by_reversal"`
}

func ValidateLedgerPlan(plan LedgerPlan) []string {
	var issues []string
	if !plan.AppendOnly {
		issues = append(issues, "append_only should be true so history is reconstructable")
	}
	if !plan.DoubleEntry {
		issues = append(issues, "double_entry should be true so every posting batch stays balanced")
	}
	if !plan.IdempotentWrites {
		issues = append(issues, "idempotent_writes should be true so retries do not duplicate money movement")
	}
	if !plan.ImmutableAuditTrail {
		issues = append(issues, "immutable_audit_trail should be true so investigators can trust history")
	}
	if !plan.ReconciliationWorkflow {
		issues = append(issues, "reconciliation_workflow should be true so external settlement drift is handled explicitly")
	}
	if !plan.ProjectionRebuild {
		issues = append(issues, "projection_rebuild should be true so balances can be regenerated from source-of-truth entries")
	}
	if !plan.CorrectionByReversal {
		issues = append(issues, "correction_by_reversal should be true so fixes do not mutate committed postings")
	}
	return issues
}

func main() {
	name := flag.String("name", "payment-ledger", "plan name")
	flag.Parse()

	plan := LedgerPlan{
		Name:                   *name,
		AppendOnly:             true,
		DoubleEntry:            true,
		IdempotentWrites:       true,
		ImmutableAuditTrail:    true,
		ReconciliationWorkflow: true,
		ProjectionRebuild:      true,
		CorrectionByReversal:   true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateLedgerPlan(plan),
	})
}
