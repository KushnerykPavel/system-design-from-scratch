package main

import "testing"

func TestValidateWalletPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  WalletPlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: WalletPlan{
				TracksAvailableAndHeld:  true,
				HasHoldExpiry:           true,
				IdempotentSettlement:    true,
				IdempotentRelease:       true,
				PreventsNegativeBalance: true,
				SupportsPartialSettle:   true,
				HasAuditTrail:           true,
			},
			wantN: 0,
		},
		{
			name: "missing reservation lifecycle fails",
			plan: WalletPlan{
				TracksAvailableAndHeld:  true,
				PreventsNegativeBalance: true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateWalletPlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
