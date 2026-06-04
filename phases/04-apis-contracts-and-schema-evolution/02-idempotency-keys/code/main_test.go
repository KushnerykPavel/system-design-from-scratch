package main

import "testing"

func TestDecide(t *testing.T) {
	tests := []struct {
		name     string
		existing *Record
		hash     string
		want     Decision
	}{
		{name: "new request executes", existing: nil, hash: "abc", want: Execute},
		{name: "matching retry replays", existing: &Record{RequestHash: "abc"}, hash: "abc", want: Replay},
		{name: "mismatched payload conflicts", existing: &Record{RequestHash: "abc"}, hash: "def", want: Conflict},
	}

	for _, tt := range tests {
		if got := Decide(tt.existing, tt.hash); got != tt.want {
			t.Fatalf("%s: got %q want %q", tt.name, got, tt.want)
		}
	}
}
