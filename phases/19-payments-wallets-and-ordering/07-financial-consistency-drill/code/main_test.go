package main

import "testing"

func TestValidateDrillScorecard(t *testing.T) {
	tests := []struct {
		name  string
		card  DrillScorecard
		wantN int
	}{
		{
			name: "complete scorecard passes",
			card: DrillScorecard{
				DefinesInvariant:        true,
				NamesSourceOfTruth:      true,
				SizesRetryAmplification: true,
				HasDeepDive:             true,
				CoversFailureModes:      true,
				CoversObservability:     true,
				HandlesRedesign:         true,
			},
			wantN: 0,
		},
		{
			name: "shallow answer fails",
			card: DrillScorecard{
				DefinesInvariant: true,
			},
			wantN: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateDrillScorecard(tt.card)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
