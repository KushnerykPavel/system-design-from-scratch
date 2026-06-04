package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
)

type RingNode struct {
	ID     string
	Weight int
	VNodes int
}

type token struct {
	Point uint64
	Node  string
}

type Ring struct {
	tokens []token
}

func buildRing(nodes []RingNode) Ring {
	var tokens []token
	for _, node := range nodes {
		totalVNodes := node.VNodes * max(node.Weight, 1)
		for i := 0; i < totalVNodes; i++ {
			tokens = append(tokens, token{
				Point: hash(fmt.Sprintf("%s#%d", node.ID, i)),
				Node:  node.ID,
			})
		}
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Point < tokens[j].Point
	})

	return Ring{tokens: tokens}
}

func (r Ring) owner(key string) string {
	if len(r.tokens) == 0 {
		return ""
	}

	point := hash(key)
	idx := sort.Search(len(r.tokens), func(i int) bool {
		return r.tokens[i].Point >= point
	})
	if idx == len(r.tokens) {
		idx = 0
	}
	return r.tokens[idx].Node
}

func remapRatio(before, after Ring, keys []string) float64 {
	if len(keys) == 0 {
		return 0
	}

	moved := 0
	for _, key := range keys {
		if before.owner(key) != after.owner(key) {
			moved++
		}
	}

	return float64(moved) / float64(len(keys))
}

func ownershipShare(r Ring, keys []string) map[string]int {
	result := make(map[string]int)
	for _, key := range keys {
		result[r.owner(key)]++
	}
	return result
}

func hash(value string) uint64 {
	sum := sha256.Sum256([]byte(value))
	return binary.BigEndian.Uint64(sum[:8])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}
	before := buildRing([]RingNode{
		{ID: "a", Weight: 1, VNodes: 32},
		{ID: "b", Weight: 1, VNodes: 32},
		{ID: "c", Weight: 1, VNodes: 32},
	})
	after := buildRing([]RingNode{
		{ID: "a", Weight: 1, VNodes: 32},
		{ID: "b", Weight: 1, VNodes: 32},
		{ID: "c", Weight: 1, VNodes: 32},
		{ID: "d", Weight: 1, VNodes: 32},
	})

	fmt.Printf("ownership before: %#v\n", ownershipShare(before, keys))
	fmt.Printf("ownership after: %#v\n", ownershipShare(after, keys))
	fmt.Printf("remap ratio: %.2f\n", remapRatio(before, after, keys))
}
