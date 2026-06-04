package main

import "testing"

func TestRecommendBucketSeconds(t *testing.T) {
	tests := []struct {
		window int
		want   int
	}{
		{window: 1800, want: 60},
		{window: 7200, want: 300},
		{window: 172800, want: 3600},
	}

	for _, test := range tests {
		if got := RecommendBucketSeconds(test.window); got != test.want {
			t.Fatalf("window=%d got=%d want=%d", test.window, got, test.want)
		}
	}
}
