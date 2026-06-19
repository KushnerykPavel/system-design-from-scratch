package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Graph represents a professional graph as an adjacency list.
// memberID → list of connected memberIDs.
type Graph map[string][]string

// AddEdge adds a bidirectional connection between members a and b.
// Duplicate edges are not deduplicated — callers should not add the same edge twice.
func AddEdge(g Graph, a, b string) {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
}

// BFS traverses the graph from start up to maxDepth hops.
// It returns a map of memberID → degree (1 = direct connection, 2 = 2nd-degree, etc.).
// The start member itself is not included in the result.
func BFS(g Graph, start string, maxDepth int) map[string]int {
	visited := make(map[string]int) // memberID → degree
	visited[start] = 0
	queue := []string{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDegree := visited[current]
		if currentDegree >= maxDepth {
			continue
		}
		for _, neighbor := range g[current] {
			if _, seen := visited[neighbor]; !seen {
				visited[neighbor] = currentDegree + 1
				queue = append(queue, neighbor)
			}
		}
	}

	// Remove the start node from results.
	delete(visited, start)
	return visited
}

// PYMK returns 2nd-degree connections for memberID that are not already directly connected.
// Results are sorted alphabetically for deterministic output.
func PYMK(g Graph, memberID string) []string {
	// Collect 1st-degree connections.
	directSet := make(map[string]bool)
	for _, peer := range g[memberID] {
		directSet[peer] = true
	}

	// BFS to depth 2.
	bfsResult := BFS(g, memberID, 2)

	// Keep only degree-2 members.
	var candidates []string
	for id, degree := range bfsResult {
		if degree == 2 && !directSet[id] {
			candidates = append(candidates, id)
		}
	}

	sort.Strings(candidates)
	return candidates
}

func main() {
	g := make(Graph)

	// Build a small professional graph.
	// alice is connected to bob, carol, dave
	AddEdge(g, "alice", "bob")
	AddEdge(g, "alice", "carol")
	AddEdge(g, "alice", "dave")

	// bob is connected to alice, eve, frank
	AddEdge(g, "bob", "eve")
	AddEdge(g, "bob", "frank")

	// carol is connected to alice, grace, hank
	AddEdge(g, "carol", "grace")
	AddEdge(g, "carol", "hank")

	// dave is connected to alice, ivan
	AddEdge(g, "dave", "ivan")

	// eve is connected to bob only (already captured)
	// frank is connected to bob, grace (cross 2nd-degree link)
	AddEdge(g, "frank", "grace")

	// Compute PYMK for alice.
	pymkAlice := PYMK(g, "alice")

	// Compute BFS degrees from alice up to depth 2.
	bfsDegrees := BFS(g, "alice", 2)

	output := map[string]any{
		"pymk_for_alice":       pymkAlice,
		"bfs_degrees_for_alice": bfsDegrees,
		"note":                 "pymk excludes direct connections (bob, carol, dave) and includes 2nd-degree only",
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
