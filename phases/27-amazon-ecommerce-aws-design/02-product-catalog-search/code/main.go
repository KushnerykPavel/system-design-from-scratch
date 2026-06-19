package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Product represents a single catalog entry.
type Product struct {
	ID            string
	Title         string
	Category      string
	Price         float64
	Rating        float64
	PrimeEligible bool
}

// SearchFilters narrows the result set before ranking.
type SearchFilters struct {
	Category  string  // empty means all categories
	MaxPrice  float64 // 0 means no price ceiling
	PrimeOnly bool
	MinRating float64 // 0 means no minimum
}

// SearchIndex is an in-memory search store.
type SearchIndex struct {
	products map[string]Product
}

// NewSearchIndex creates an empty index.
func NewSearchIndex() *SearchIndex {
	return &SearchIndex{products: make(map[string]Product)}
}

// Index adds or replaces a product in the index.
func (idx *SearchIndex) Index(p Product) {
	idx.products[p.ID] = p
}

// Search returns products matching the query and filters, ranked by relevance.
func (idx *SearchIndex) Search(query string, filters SearchFilters) []Product {
	query = strings.ToLower(strings.TrimSpace(query))
	var candidates []Product
	for _, p := range idx.products {
		if !matchesFilters(p, filters) {
			continue
		}
		if query == "" || matchesQuery(p, query) {
			candidates = append(candidates, p)
		}
	}
	return RankProducts(candidates, query)
}

// matchesFilters returns true when the product satisfies all active filters.
func matchesFilters(p Product, f SearchFilters) bool {
	if f.Category != "" && !strings.EqualFold(p.Category, f.Category) {
		return false
	}
	if f.MaxPrice > 0 && p.Price > f.MaxPrice {
		return false
	}
	if f.PrimeOnly && !p.PrimeEligible {
		return false
	}
	if f.MinRating > 0 && p.Rating < f.MinRating {
		return false
	}
	return true
}

// matchesQuery returns true when title or category contains the query terms.
func matchesQuery(p Product, query string) bool {
	haystack := strings.ToLower(p.Title + " " + p.Category)
	for _, term := range strings.Fields(query) {
		if !strings.Contains(haystack, term) {
			return false
		}
	}
	return true
}

// relevanceScore computes a simple ranking score for a product given a query.
// Higher is better.
func relevanceScore(p Product, query string) float64 {
	score := 0.0

	// Text relevance: count how many query terms appear in the title.
	titleLower := strings.ToLower(p.Title)
	for _, term := range strings.Fields(query) {
		if strings.Contains(titleLower, term) {
			score += 3.0 // title match is weighted higher
		}
	}

	// Rating boost: normalise 0–5 scale to 0–2 points.
	score += (p.Rating / 5.0) * 2.0

	// Prime eligibility bonus.
	if p.PrimeEligible {
		score += 1.0
	}

	return score
}

// RankProducts sorts products by relevance score descending.
func RankProducts(products []Product, query string) []Product {
	ranked := make([]Product, len(products))
	copy(ranked, products)
	sort.SliceStable(ranked, func(i, j int) bool {
		si := relevanceScore(ranked[i], query)
		sj := relevanceScore(ranked[j], query)
		if si != sj {
			return si > sj
		}
		// Tie-break: higher rating first, then alphabetical title.
		if ranked[i].Rating != ranked[j].Rating {
			return ranked[i].Rating > ranked[j].Rating
		}
		return ranked[i].Title < ranked[j].Title
	})
	return ranked
}

func main() {
	idx := NewSearchIndex()

	// Seed the index with sample products.
	products := []Product{
		{ID: "p1", Title: "Wireless Bluetooth Headphones", Category: "Electronics", Price: 79.99, Rating: 4.5, PrimeEligible: true},
		{ID: "p2", Title: "Noise Cancelling Headphones Pro", Category: "Electronics", Price: 249.99, Rating: 4.8, PrimeEligible: true},
		{ID: "p3", Title: "Budget Wired Headphones", Category: "Electronics", Price: 19.99, Rating: 3.2, PrimeEligible: false},
		{ID: "p4", Title: "Running Shoes Ultra Comfort", Category: "Sports", Price: 129.99, Rating: 4.6, PrimeEligible: true},
		{ID: "p5", Title: "Yoga Mat Premium Non-Slip", Category: "Sports", Price: 39.99, Rating: 4.1, PrimeEligible: false},
	}
	for _, p := range products {
		idx.Index(p)
	}

	// Demo 1: search "headphones" with Prime filter and max price.
	fmt.Println("=== Search: 'headphones', Prime only, max $100 ===")
	results := idx.Search("headphones", SearchFilters{PrimeOnly: true, MaxPrice: 100})
	printResults(results)

	// Demo 2: search all sports products rated >= 4.5.
	fmt.Println("\n=== Search: 'sports', min rating 4.5 ===")
	results2 := idx.Search("", SearchFilters{Category: "Sports", MinRating: 4.5})
	printResults(results2)

	// Demo 3: no query, no filters — return all ranked by score.
	fmt.Println("\n=== Search: all products ranked ===")
	results3 := idx.Search("", SearchFilters{})
	printResults(results3)
}

func printResults(products []Product) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(products)
}
