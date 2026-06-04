package main

import "testing"

func TestBottleneck(t *testing.T) {
	name, limit := Bottleneck(BottleneckModel{
		CPUReqPerSecond:     128000,
		DiskReqPerSecond:    20000,
		NetworkReqPerSecond: 131000,
	})

	if name != "disk" {
		t.Fatalf("unexpected bottleneck: %s", name)
	}
	if limit != 20000 {
		t.Fatalf("unexpected limit: %.2f", limit)
	}
}
