package main

import "testing"

func TestPlanTransaction(t *testing.T) {
	tests := []struct {
		name     string
		profile  TransactionProfile
		want     IsolationLevel
		wantSaga bool
		warnings int
	}{
		{
			name:    "simple noncritical write",
			profile: TransactionProfile{},
			want:    ReadCommitted,
		},
		{
			name: "critical invariant high conflict",
			profile: TransactionProfile{
				InvariantCritical: true,
				ConflictRate:      "high",
			},
			want: Serializable,
		},
		{
			name: "cross service hotspot",
			profile: TransactionProfile{
				InvariantCritical: true,
				ConflictRate:      "high",
				HotspotRisk:       "high",
				CrossService:      true,
			},
			want:     Serializable,
			wantSaga: true,
			warnings: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlanTransaction(tt.profile)
			if got.Level != tt.want {
				t.Fatalf("PlanTransaction() level = %q, want %q", got.Level, tt.want)
			}
			if got.UseSaga != tt.wantSaga {
				t.Fatalf("PlanTransaction() saga = %v, want %v", got.UseSaga, tt.wantSaga)
			}
			if len(got.Warnings) != tt.warnings {
				t.Fatalf("PlanTransaction() warnings = %d, want %d", len(got.Warnings), tt.warnings)
			}
		})
	}
}
