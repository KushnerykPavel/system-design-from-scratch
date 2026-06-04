package main

import "testing"

func TestScoreCapacityRound(t *testing.T) {
	tests := []struct {
		name  string
		round CapacityRound
		wantN int
	}{
		{
			name: "complete round passes",
			round: CapacityRound{
				HasQPS:             true,
				HasPeakFactor:      true,
				HasStorageOrEgress: true,
				HasAmplification:   true,
				NamesBottleneck:    true,
				LinksToDesign:      true,
			},
			wantN: 0,
		},
		{
			name: "partial round fails",
			round: CapacityRound{
				HasQPS: true,
			},
			wantN: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreCapacityRound(tt.round)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
