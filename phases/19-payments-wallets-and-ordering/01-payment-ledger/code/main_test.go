package main

import "testing"

func TestValidateLedgerPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  LedgerPlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: LedgerPlan{
				AppendOnly:             true,
				DoubleEntry:            true,
				IdempotentWrites:       true,
				ImmutableAuditTrail:    true,
				ReconciliationWorkflow: true,
				ProjectionRebuild:      true,
				CorrectionByReversal:   true,
			},
			wantN: 0,
		},
		{
			name: "missing financial controls fails",
			plan: LedgerPlan{
				AppendOnly:       true,
				IdempotentWrites: true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateLedgerPlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
