package main

import "testing"

func TestValidateCompliancePlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  CompliancePlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: CompliancePlan{
				ImmutableAudit:     true,
				PIISeparation:      true,
				DeletionWorkflow:   true,
				LegalHoldSupport:   true,
				AccessLogging:      true,
				PolicyByRecordType: true,
				ArchiveStrategy:    true,
			},
			wantN: 0,
		},
		{
			name: "missing governance controls fails",
			plan: CompliancePlan{
				ImmutableAudit: true,
				AccessLogging:  true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCompliancePlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
