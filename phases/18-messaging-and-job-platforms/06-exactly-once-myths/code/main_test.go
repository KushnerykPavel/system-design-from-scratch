package main

import "testing"

func TestValidateExactlyOnceClaim(t *testing.T) {
	tests := []struct {
		name  string
		claim ExactlyOnceClaim
		wantN int
	}{
		{
			name: "bounded claim",
			claim: ExactlyOnceClaim{
				DefinesBoundary:       true,
				HasIdempotentConsumer: true,
				HasDedupKey:           true,
				ExplainsFailureCase:   true,
				NamesResidualRisk:     true,
			},
			wantN: 0,
		},
		{
			name: "unsafe side-effect claim",
			claim: ExactlyOnceClaim{
				DefinesBoundary:       true,
				HasIdempotentConsumer: true,
				HasDedupKey:           true,
				IncludesSideEffects:   true,
				ExplainsFailureCase:   true,
			},
			wantN: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(ValidateExactlyOnceClaim(tt.claim)); got != tt.wantN {
				t.Fatalf("got %d issues, want %d", got, tt.wantN)
			}
		})
	}
}
