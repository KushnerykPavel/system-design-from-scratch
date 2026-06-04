package main

import "testing"

func TestValidateRubricScorecard(t *testing.T) {
	tests := []struct {
		name  string
		card  RubricScorecard
		wantN int
	}{
		{
			name: "complete answer passes",
			card: RubricScorecard{
				ClarifiesScope:         true,
				QuantifiesWorkload:     true,
				HasHighLevelDesign:     true,
				ChoosesDeepDive:        true,
				ExplainsTradeoffs:      true,
				CoversFailureModes:     true,
				CoversObservability:    true,
				HandlesRedesignCleanly: true,
			},
			wantN: 0,
		},
		{
			name: "shallow answer fails",
			card: RubricScorecard{
				ClarifiesScope: true,
			},
			wantN: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateRubricScorecard(tt.card)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
