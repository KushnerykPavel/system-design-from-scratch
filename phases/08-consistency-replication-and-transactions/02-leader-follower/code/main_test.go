package main

import "testing"

func TestAssess(t *testing.T) {
	tests := []struct {
		name        string
		topology    Topology
		wantSafe    bool
		wantRisk    string
		wantWarning int
	}{
		{
			name: "stronger topology",
			topology: Topology{
				Followers:         3,
				AckPolicy:         AckMajority,
				FencingEnabled:    true,
				CandidateCaughtUp: true,
			},
			wantSafe: true,
			wantRisk: "low",
		},
		{
			name: "critical follower reads with lag",
			topology: Topology{
				Followers:             2,
				AckPolicy:             AckLeaderPlus1,
				CriticalFollowerReads: true,
				MaxFollowerLagSeconds: 2,
				FencingEnabled:        true,
				CandidateCaughtUp:     true,
			},
			wantSafe:    false,
			wantRisk:    "moderate",
			wantWarning: 1,
		},
		{
			name: "leader local ack and unfenced auto failover",
			topology: Topology{
				Followers:         2,
				AckPolicy:         AckLeaderLocal,
				AutoFailover:      true,
				FencingEnabled:    false,
				CandidateCaughtUp: false,
			},
			wantSafe:    true,
			wantRisk:    "high",
			wantWarning: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Assess(tt.topology)
			if got.SafeForCriticalReads != tt.wantSafe {
				t.Fatalf("Assess() safe = %v, want %v", got.SafeForCriticalReads, tt.wantSafe)
			}
			if got.FailoverRisk != tt.wantRisk {
				t.Fatalf("Assess() risk = %q, want %q", got.FailoverRisk, tt.wantRisk)
			}
			if len(got.Warnings) != tt.wantWarning {
				t.Fatalf("Assess() warnings = %d, want %d", len(got.Warnings), tt.wantWarning)
			}
		})
	}
}
