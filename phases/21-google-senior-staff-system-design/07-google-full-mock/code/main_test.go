package main

import "testing"

func TestValidateMockScorecard(t *testing.T) {
	tests := []struct {
		name  string
		card  MockScorecard
		wantN int
	}{
		{
			name: "complete mock passes",
			card: MockScorecard{
				ClarifiesPrompt:         true,
				PrioritizesRequirements: true,
				IncludesSizing:          true,
				HasHighLevelDesign:      true,
				HasDeepDive:             true,
				CoversRiskAndOps:        true,
				HandlesRedesign:         true,
				StaysTimeBound:          true,
			},
			wantN: 0,
		},
		{
			name: "partial mock fails",
			card: MockScorecard{
				ClarifiesPrompt: true,
			},
			wantN: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateMockScorecard(tt.card)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
