package main

import (
	"fmt"
	"testing"
)

func TestConsistentHashingRemapsMinorityOfKeysOnSingleAdd(t *testing.T) {
	var keys []string
	for i := 0; i < 1000; i++ {
		keys = append(keys, fmt.Sprintf("key-%d", i))
	}

	before := buildRing([]RingNode{
		{ID: "a", Weight: 1, VNodes: 64},
		{ID: "b", Weight: 1, VNodes: 64},
		{ID: "c", Weight: 1, VNodes: 64},
	})
	after := buildRing([]RingNode{
		{ID: "a", Weight: 1, VNodes: 64},
		{ID: "b", Weight: 1, VNodes: 64},
		{ID: "c", Weight: 1, VNodes: 64},
		{ID: "d", Weight: 1, VNodes: 64},
	})

	ratio := remapRatio(before, after, keys)
	if ratio <= 0 || ratio >= 0.5 {
		t.Fatalf("expected bounded remap ratio, got %.3f", ratio)
	}
}

func TestWeightedPlacementAssignsMoreKeysToHeavierNode(t *testing.T) {
	var keys []string
	for i := 0; i < 2000; i++ {
		keys = append(keys, fmt.Sprintf("key-%d", i))
	}

	ring := buildRing([]RingNode{
		{ID: "small", Weight: 1, VNodes: 32},
		{ID: "large", Weight: 2, VNodes: 32},
	})

	share := ownershipShare(ring, keys)
	if share["large"] <= share["small"] {
		t.Fatalf("expected weighted node to own more keys, got %#v", share)
	}
}
