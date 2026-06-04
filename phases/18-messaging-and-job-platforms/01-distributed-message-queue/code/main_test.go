package main

import "testing"

func TestValidateQueuePlan(t *testing.T) {
	tests := []struct {
		name  string
		plan  QueuePlan
		wantN int
	}{
		{
			name: "complete plan passes",
			plan: QueuePlan{
				HasDurableLog:            true,
				UsesVisibilityTimeout:    true,
				TracksConsumerOffsets:    true,
				SupportsRedelivery:       true,
				HasPartitionStrategy:     true,
				HasPoisonMessageHandling: true,
				HasBackpressureControls:  true,
			},
			wantN: 0,
		},
		{
			name: "missing critical controls returns issues",
			plan: QueuePlan{
				HasDurableLog:         true,
				TracksConsumerOffsets: true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateQueuePlan(tt.plan)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
