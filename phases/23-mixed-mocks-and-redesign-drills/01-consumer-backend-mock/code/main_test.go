package main

import "testing"

func TestValidateConsumerMock(t *testing.T) {
	tests := []struct {
		name  string
		card  ConsumerMockScorecard
		wantN int
	}{
		{
			name: "complete mock passes",
			card: ConsumerMockScorecard{
				ClarifiesUserPromise:    true,
				SizesReadWriteAsymmetry: true,
				ChoosesStateBoundary:    true,
				CoversSkewOrFanout:      true,
				NamesDegradedMode:       true,
				CoversObservability:     true,
				HandlesRedesign:         true,
			},
			wantN: 0,
		},
		{
			name: "partial mock fails",
			card: ConsumerMockScorecard{
				ClarifiesUserPromise: true,
			},
			wantN: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateConsumerMock(tt.card)
			if len(got) != tt.wantN {
				t.Fatalf("got %d issues, want %d: %v", len(got), tt.wantN, got)
			}
		})
	}
}
