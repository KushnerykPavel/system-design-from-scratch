package main

import (
	"math"
	"testing"
)

func TestCosineSimilarityIdentical(t *testing.T) {
	v := []float64{1.0, 0.5, 0.25}
	score := cosineSimilarity(v, v)
	if math.Abs(score-1.0) > 1e-9 {
		t.Fatalf("expected cosine similarity of 1.0 for identical vectors, got %f", score)
	}
}

func TestCosineSimilarityOrthogonal(t *testing.T) {
	a := []float64{1.0, 0.0}
	b := []float64{0.0, 1.0}
	score := cosineSimilarity(a, b)
	if math.Abs(score) > 1e-9 {
		t.Fatalf("expected 0.0 for orthogonal vectors, got %f", score)
	}
}

func TestRetrieveTopK(t *testing.T) {
	items := []ItemEmbedding{
		{ItemID: "a", Embedding: []float64{1.0, 0.0}},
		{ItemID: "b", Embedding: []float64{0.0, 1.0}},
		{ItemID: "c", Embedding: []float64{0.7, 0.3}},
	}
	user := []float64{1.0, 0.0}
	top2 := Retrieve(user, items, 2)
	if len(top2) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top2))
	}
	if top2[0].ItemID != "a" {
		t.Fatalf("expected top result to be 'a', got '%s'", top2[0].ItemID)
	}
}

func TestRetrieveKGreaterThanItems(t *testing.T) {
	items := []ItemEmbedding{
		{ItemID: "x", Embedding: []float64{1.0}},
	}
	results := Retrieve([]float64{1.0}, items, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result when K > item count, got %d", len(results))
	}
}
