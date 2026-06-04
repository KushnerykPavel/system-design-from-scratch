package main

import "testing"

func TestValidateConsistencyScenario(t *testing.T) {
	tests := []struct {
		name  string
		input ConsistencyScenario
		wantN int
	}{
		{
			name: "complete scenario passes",
			input: ConsistencyScenario{
				DefinesSourceOfTruth:     true,
				NamesStrongPath:          true,
				NamesStaleReadAllowance:  true,
				DefinesAnomalyBudget:     true,
				ExplainsFailoverBehavior: true,
				CoversLagMetrics:         true,
				StatesTradeoffs:          true,
			},
			wantN: 0,
		},
		{
			name: "underspecified scenario fails",
			input: ConsistencyScenario{
				DefinesSourceOfTruth: true,
			},
			wantN: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateConsistencyScenario(tt.input)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
