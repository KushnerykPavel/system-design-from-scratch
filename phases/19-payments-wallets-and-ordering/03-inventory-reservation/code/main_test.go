package main

import "testing"

func TestValidateReservationPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  ReservationPlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: ReservationPlan{
				AuthoritativeReservePath: true,
				ReservationTTL:           true,
				IdempotentConfirm:        true,
				IdempotentRelease:        true,
				OversellGuard:            true,
				HotSKUProtection:         true,
				LeakDetection:            true,
			},
			wantN: 0,
		},
		{
			name: "missing stock protections fails",
			plan: ReservationPlan{
				ReservationTTL: true,
				LeakDetection:  true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateReservationPlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
