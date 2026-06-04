package main

import "testing"

func TestValidateInfraPlatform(t *testing.T) {
	tests := []struct {
		name  string
		card  InfraPlatformScorecard
		wantN int
	}{
		{
			name: "complete mock passes",
			card: InfraPlatformScorecard{
				ClarifiesPlatformContract:    true,
				SeparatesControlAndDataPlane: true,
				SizesHotAndSlowPaths:         true,
				DefinesFailureDomain:         true,
				HasRolloutOrMigrationPlan:    true,
				CoversTenantGuardrails:       true,
				CoversObservability:          true,
				HandlesConstraintChange:      true,
			},
			wantN: 0,
		},
		{
			name: "partial mock fails",
			card: InfraPlatformScorecard{
				ClarifiesPlatformContract: true,
			},
			wantN: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateInfraPlatform(tt.card)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
