package main

import (
	"encoding/json"
	"flag"
	"os"
)

type CompliancePlan struct {
	Name               string `json:"name"`
	ImmutableAudit     bool   `json:"immutable_audit"`
	PIISeparation      bool   `json:"pii_separation"`
	DeletionWorkflow   bool   `json:"deletion_workflow"`
	LegalHoldSupport   bool   `json:"legal_hold_support"`
	AccessLogging      bool   `json:"access_logging"`
	PolicyByRecordType bool   `json:"policy_by_record_type"`
	ArchiveStrategy    bool   `json:"archive_strategy"`
}

func ValidateCompliancePlan(plan CompliancePlan) []string {
	var issues []string
	if !plan.ImmutableAudit {
		issues = append(issues, "immutable_audit should be true so privileged actions remain trustworthy")
	}
	if !plan.PIISeparation {
		issues = append(issues, "pii_separation should be true so privacy and PCI blast radius stay smaller")
	}
	if !plan.DeletionWorkflow {
		issues = append(issues, "deletion_workflow should be true so privacy requests are handled systematically")
	}
	if !plan.LegalHoldSupport {
		issues = append(issues, "legal_hold_support should be true so deletion cannot violate legal constraints")
	}
	if !plan.AccessLogging {
		issues = append(issues, "access_logging should be true so reads and exports of sensitive data are traceable")
	}
	if !plan.PolicyByRecordType {
		issues = append(issues, "policy_by_record_type should be true so retention is driven by record class, not one global rule")
	}
	if !plan.ArchiveStrategy {
		issues = append(issues, "archive_strategy should be true so long retention does not overload hot storage")
	}
	return issues
}

func main() {
	name := flag.String("name", "audit-and-compliance", "plan name")
	flag.Parse()

	plan := CompliancePlan{
		Name:               *name,
		ImmutableAudit:     true,
		PIISeparation:      true,
		DeletionWorkflow:   true,
		LegalHoldSupport:   true,
		AccessLogging:      true,
		PolicyByRecordType: true,
		ArchiveStrategy:    true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateCompliancePlan(plan),
	})
}
