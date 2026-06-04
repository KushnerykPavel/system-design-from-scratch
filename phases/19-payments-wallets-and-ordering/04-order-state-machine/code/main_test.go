package main

import "testing"

func TestValidateRecoveryPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  RecoveryPlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: RecoveryPlan{
				ExplicitStates:        true,
				IdempotentTransitions: true,
				TimeoutStates:         true,
				CompensationPath:      true,
				EventHistory:          true,
				OperatorRecovery:      true,
				CallbackDedupe:        true,
			},
			wantN: 0,
		},
		{
			name: "missing recovery controls fails",
			plan: RecoveryPlan{
				ExplicitStates: true,
				EventHistory:   true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateRecoveryPlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
