package main

import "testing"

func TestValidateSchedulerPlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  SchedulerPlan
		wantN int
	}{
		{
			name: "healthy scheduler",
			plan: SchedulerPlan{
				HasShardLeases:       true,
				HasRetryPolicy:       true,
				HasJitter:            true,
				HandlesMissedRuns:    true,
				HasTenantQuotas:      true,
				HasDLQ:               true,
				HasDeadlineAwareness: true,
			},
			wantN: 0,
		},
		{
			name: "weak scheduler",
			plan: SchedulerPlan{
				HasRetryPolicy: true,
			},
			wantN: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(ValidateSchedulerPlan(tt.plan)); got != tt.wantN {
				t.Fatalf("got %d issues, want %d", got, tt.wantN)
			}
		})
	}
}
