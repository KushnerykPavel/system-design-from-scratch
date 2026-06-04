package main

import "testing"

func TestChooseInterface(t *testing.T) {
	tests := []struct {
		name string
		in   Workload
		want InterfaceKind
	}{
		{
			name: "public client request stays http",
			in: Workload{
				ExternalClients:     true,
				NeedsImmediateReply: true,
			},
			want: HTTP,
		},
		{
			name: "internal low latency call prefers grpc",
			in: Workload{
				NeedsImmediateReply: true,
				LowLatencyInternal:  true,
			},
			want: GRPC,
		},
		{
			name: "fanout workflow becomes event",
			in: Workload{
				HighFanout: true,
			},
			want: Event,
		},
	}

	for _, tt := range tests {
		if got := ChooseInterface(tt.in); got != tt.want {
			t.Fatalf("%s: got %q want %q", tt.name, got, tt.want)
		}
	}
}
