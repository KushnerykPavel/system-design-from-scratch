package main

import "testing"

func seedIndex() *SearchIndex {
	idx := NewSearchIndex()
	idx.Index(Product{ID: "p1", Title: "Wireless Bluetooth Headphones", Category: "Electronics", Price: 79.99, Rating: 4.5, PrimeEligible: true})
	idx.Index(Product{ID: "p2", Title: "Noise Cancelling Headphones Pro", Category: "Electronics", Price: 249.99, Rating: 4.8, PrimeEligible: true})
	idx.Index(Product{ID: "p3", Title: "Budget Wired Headphones", Category: "Electronics", Price: 19.99, Rating: 3.2, PrimeEligible: false})
	idx.Index(Product{ID: "p4", Title: "Running Shoes Ultra Comfort", Category: "Sports", Price: 129.99, Rating: 4.6, PrimeEligible: true})
	idx.Index(Product{ID: "p5", Title: "Yoga Mat Premium Non-Slip", Category: "Sports", Price: 39.99, Rating: 4.1, PrimeEligible: false})
	return idx
}

func TestIndexAndRetrieveAll(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("", SearchFilters{})
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}

func TestFilterByCategory(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("", SearchFilters{Category: "Sports"})
	if len(results) != 2 {
		t.Fatalf("expected 2 Sports products, got %d", len(results))
	}
	for _, p := range results {
		if p.Category != "Sports" {
			t.Errorf("unexpected category %q in results", p.Category)
		}
	}
}

func TestFilterPrimeOnly(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("", SearchFilters{PrimeOnly: true})
	for _, p := range results {
		if !p.PrimeEligible {
			t.Errorf("non-Prime product %q returned with PrimeOnly filter", p.ID)
		}
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 Prime products, got %d", len(results))
	}
}

func TestFilterMaxPrice(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("", SearchFilters{MaxPrice: 50.0})
	for _, p := range results {
		if p.Price > 50.0 {
			t.Errorf("product %q price %.2f exceeds max price 50.0", p.ID, p.Price)
		}
	}
}

func TestFilterMinRating(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("", SearchFilters{MinRating: 4.5})
	for _, p := range results {
		if p.Rating < 4.5 {
			t.Errorf("product %q rating %.1f below min rating 4.5", p.ID, p.Rating)
		}
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 products with rating >= 4.5, got %d", len(results))
	}
}

func TestQueryTextMatch(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("headphones", SearchFilters{})
	if len(results) != 3 {
		t.Fatalf("expected 3 headphone products, got %d", len(results))
	}
}

func TestQueryNoMatch(t *testing.T) {
	idx := seedIndex()
	results := idx.Search("laptop", SearchFilters{})
	if len(results) != 0 {
		t.Fatalf("expected 0 results for 'laptop', got %d", len(results))
	}
}

func TestRankingHigherRatedFirst(t *testing.T) {
	// Two products with equal query relevance; higher rating must rank first.
	idx := NewSearchIndex()
	idx.Index(Product{ID: "low", Title: "Widget A", Category: "Tools", Price: 10, Rating: 3.0, PrimeEligible: false})
	idx.Index(Product{ID: "high", Title: "Widget B", Category: "Tools", Price: 10, Rating: 4.9, PrimeEligible: false})
	results := idx.Search("widget", SearchFilters{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].ID != "high" {
		t.Errorf("expected high-rated product first, got %q", results[0].ID)
	}
}

func TestRankingPrimeBonusApplied(t *testing.T) {
	// Prime product with slightly lower rating should beat non-Prime with same title match.
	idx := NewSearchIndex()
	idx.Index(Product{ID: "prime", Title: "Camera DSLR", Category: "Electronics", Price: 500, Rating: 4.0, PrimeEligible: true})
	idx.Index(Product{ID: "nonprime", Title: "Camera DSLR", Category: "Electronics", Price: 500, Rating: 4.0, PrimeEligible: false})
	results := idx.Search("camera", SearchFilters{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].ID != "prime" {
		t.Errorf("expected Prime product to rank first, got %q", results[0].ID)
	}
}

func TestIndexOverwrite(t *testing.T) {
	idx := NewSearchIndex()
	idx.Index(Product{ID: "x1", Title: "Old Title", Category: "Books", Price: 5.0, Rating: 2.0, PrimeEligible: false})
	idx.Index(Product{ID: "x1", Title: "New Title", Category: "Books", Price: 5.0, Rating: 4.9, PrimeEligible: true})
	results := idx.Search("new title", SearchFilters{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result after overwrite, got %d", len(results))
	}
	if results[0].Title != "New Title" {
		t.Errorf("expected updated title, got %q", results[0].Title)
	}
}
