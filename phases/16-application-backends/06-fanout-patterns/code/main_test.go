package main

import "testing"

func TestRecommendFanout(t *testing.T) {
	tests := []struct {
		name  string
		shape WorkloadShape
		want  string
	}{
		{
			name: "small recipient push",
			shape: WorkloadShape{
				RecipientsPerEvent: 10,
			},
			want: "push",
		},
		{
			name: "celebrity skew mixed",
			shape: WorkloadShape{
				RecipientsPerEvent: 500000,
				SkewedAudience:     true,
			},
			want: "mixed",
		},
		{
			name: "write heavy large audience pull",
			shape: WorkloadShape{
				ReadsPerWrite:      0,
				RecipientsPerEvent: 10000,
			},
			want: "pull",
		},
	}

	for _, tt := range tests {
		if got := RecommendFanout(tt.shape); got != tt.want {
			t.Fatalf("%s: RecommendFanout() = %q, want %q", tt.name, got, tt.want)
		}
	}
}
