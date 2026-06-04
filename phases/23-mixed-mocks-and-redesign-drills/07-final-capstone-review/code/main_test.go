package main

import "testing"

func TestValidatePanelReview(t *testing.T) {
	tests := []struct {
		name   string
		review PanelReview
		wantN  int
	}{
		{
			name: "valid review passes",
			review: PanelReview{
				Clarification: 4,
				Sizing:        3,
				Architecture:  4,
				DeepDive:      3,
				FailureModes:  4,
				Observability: 3,
				TradeOffs:     4,
				Communication: 4,
			},
			wantN: 0,
		},
		{
			name: "invalid scores fail",
			review: PanelReview{
				Clarification: 0,
				Sizing:        5,
				Architecture:  3,
				DeepDive:      3,
				FailureModes:  3,
				Observability: 3,
				TradeOffs:     3,
				Communication: 3,
			},
			wantN: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePanelReview(tt.review)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
