package main

import "testing"

func TestValidateFraudHookPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  FraudHookPlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: FraudHookPlan{
				InlineBoundedChecks: true,
				AsyncScoring:        true,
				ManualReviewLane:    true,
				ExplicitFallback:    true,
				ModelVersionLogging: true,
				PolicyRollback:      true,
				OutcomeAuditability: true,
			},
			wantN: 0,
		},
		{
			name: "missing degradation controls fails",
			plan: FraudHookPlan{
				AsyncScoring: true,
			},
			wantN: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFraudHookPlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
