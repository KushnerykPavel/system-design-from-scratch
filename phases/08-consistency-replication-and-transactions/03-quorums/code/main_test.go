package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name        string
		plan        Plan
		wantOverlap bool
		wantAvail   string
		wantWarn    int
	}{
		{
			name: "overlapping quorums",
			plan: Plan{
				Replicas:    3,
				ReadQuorum:  2,
				WriteQuorum: 2,
			},
			wantOverlap: true,
			wantAvail:   "balanced",
		},
		{
			name: "tiny quorums",
			plan: Plan{
				Replicas:    3,
				ReadQuorum:  1,
				WriteQuorum: 1,
			},
			wantOverlap: false,
			wantAvail:   "very high",
			wantWarn:    1,
		},
		{
			name: "read bias mismatch and conflict risk",
			plan: Plan{
				Replicas:     5,
				ReadQuorum:   3,
				WriteQuorum:  2,
				LatencyBias:  "read",
				ConflictRisk: "high",
			},
			wantOverlap: false,
			wantAvail:   "high",
			wantWarn:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(tt.plan)
			if got.Overlaps != tt.wantOverlap {
				t.Fatalf("Evaluate() overlap = %v, want %v", got.Overlaps, tt.wantOverlap)
			}
			if got.Availability != tt.wantAvail {
				t.Fatalf("Evaluate() availability = %q, want %q", got.Availability, tt.wantAvail)
			}
			if len(got.Warnings) != tt.wantWarn {
				t.Fatalf("Evaluate() warnings = %d, want %d", len(got.Warnings), tt.wantWarn)
			}
		})
	}
}
