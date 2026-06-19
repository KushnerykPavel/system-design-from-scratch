package main

import "testing"

func TestBFSFirstDegree(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "alice", "bob")
	AddEdge(g, "alice", "carol")

	result := BFS(g, "alice", 2)

	if result["bob"] != 1 {
		t.Fatalf("expected bob at degree 1, got %d", result["bob"])
	}
	if result["carol"] != 1 {
		t.Fatalf("expected carol at degree 1, got %d", result["carol"])
	}
	if _, ok := result["alice"]; ok {
		t.Fatal("start node alice should not appear in BFS result")
	}
}

func TestBFSSecondDegree(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "alice", "bob")
	AddEdge(g, "bob", "carol")

	result := BFS(g, "alice", 2)

	if result["bob"] != 1 {
		t.Fatalf("expected bob at degree 1, got %d", result["bob"])
	}
	if result["carol"] != 2 {
		t.Fatalf("expected carol at degree 2, got %d", result["carol"])
	}
}

func TestBFSDepthLimit(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "a", "b")
	AddEdge(g, "b", "c")
	AddEdge(g, "c", "d") // d is at depth 3 from a

	result := BFS(g, "a", 2)

	if _, ok := result["d"]; ok {
		t.Fatal("d at depth 3 should not appear in BFS with maxDepth=2")
	}
	if result["c"] != 2 {
		t.Fatalf("expected c at degree 2, got %d", result["c"])
	}
}

func TestPYMKExcludesDirectConnections(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "alice", "bob")
	AddEdge(g, "bob", "carol") // carol is 2nd-degree from alice
	AddEdge(g, "alice", "carol") // but also directly connected

	candidates := PYMK(g, "alice")

	// carol should NOT appear because alice is already directly connected
	for _, c := range candidates {
		if c == "carol" {
			t.Fatal("carol should be excluded from PYMK because alice is already connected to carol")
		}
	}
}

func TestPYMKReturnsTwoDegreeNotDirect(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "alice", "bob")
	AddEdge(g, "alice", "carol")
	AddEdge(g, "bob", "dave")   // dave is 2nd-degree via bob
	AddEdge(g, "carol", "eve")  // eve is 2nd-degree via carol

	candidates := PYMK(g, "alice")

	daveFound := false
	eveFound := false
	for _, c := range candidates {
		if c == "dave" {
			daveFound = true
		}
		if c == "eve" {
			eveFound = true
		}
		// bob and carol are direct connections and must not appear
		if c == "bob" || c == "carol" {
			t.Fatalf("direct connection %s should not appear in PYMK", c)
		}
	}
	if !daveFound {
		t.Fatal("expected dave (2nd-degree via bob) in PYMK candidates")
	}
	if !eveFound {
		t.Fatal("expected eve (2nd-degree via carol) in PYMK candidates")
	}
}

func TestBFSDisconnectedNode(t *testing.T) {
	g := make(Graph)
	AddEdge(g, "alice", "bob")
	// carol is not connected to anyone

	result := BFS(g, "alice", 2)

	if _, ok := result["carol"]; ok {
		t.Fatal("carol is disconnected; should not appear in BFS")
	}
}
