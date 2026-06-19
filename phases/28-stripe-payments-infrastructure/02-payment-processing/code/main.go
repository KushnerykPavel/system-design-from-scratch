package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// CachedResponse holds the stored result for an idempotency key.
type CachedResponse struct {
	RequestHash  string
	ResponseCode int
	ResponseBody string
	CreatedAt    time.Time
}

// IdempotencyStore is an in-memory deduplication store keyed by idempotency key.
type IdempotencyStore struct {
	entries map[string]CachedResponse
}

// NewIdempotencyStore creates an empty store.
func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{entries: make(map[string]CachedResponse)}
}

// hashBody returns a hex SHA-256 digest of the request body.
func hashBody(body string) string {
	sum := sha256.Sum256([]byte(body))
	return fmt.Sprintf("%x", sum)
}

// Store records a new idempotency entry. Returns true if the entry was stored
// (key was new), false if the key already exists (duplicate, not stored again).
func (s *IdempotencyStore) Store(key, reqBody string, code int, responseBody string) bool {
	if _, exists := s.entries[key]; exists {
		return false
	}
	s.entries[key] = CachedResponse{
		RequestHash:  hashBody(reqBody),
		ResponseCode: code,
		ResponseBody: responseBody,
		CreatedAt:    time.Now(),
	}
	return true
}

// ErrHashMismatch is returned when a retry supplies a different request body
// for an already-stored idempotency key.
var ErrHashMismatch = errors.New("idempotency key reused with different request body")

// Lookup retrieves a cached response for the given key and request body.
// Returns (entry, true, nil) on cache hit with matching hash.
// Returns (zero, false, nil) when the key is not found.
// Returns (zero, false, ErrHashMismatch) when the key exists but the body hash differs.
func (s *IdempotencyStore) Lookup(key, reqBody string) (CachedResponse, bool, error) {
	entry, exists := s.entries[key]
	if !exists {
		return CachedResponse{}, false, nil
	}
	if entry.RequestHash != hashBody(reqBody) {
		return CachedResponse{}, false, ErrHashMismatch
	}
	return entry, true, nil
}

// IsExpired reports whether a cached entry is older than windowHours.
func IsExpired(entry CachedResponse, windowHours int) bool {
	return time.Since(entry.CreatedAt) > time.Duration(windowHours)*time.Hour
}

func main() {
	store := NewIdempotencyStore()

	idempotencyKey := "idem-key-abc123"
	requestBody := `{"amount":5000,"currency":"usd","payment_method":"pm_card_visa"}`

	// First attempt: store the result.
	stored := store.Store(idempotencyKey, requestBody, 200, `{"id":"pi_1","status":"succeeded"}`)
	fmt.Printf("First attempt stored: %v\n", stored)

	// Retry with same key and same body: should return cache hit.
	entry, hit, err := store.Lookup(idempotencyKey, requestBody)
	if err != nil {
		fmt.Printf("Retry error: %v\n", err)
	} else if hit {
		fmt.Printf("Cache hit — returning stored response: code=%d body=%s\n",
			entry.ResponseCode, entry.ResponseBody)
	}

	// Tampered retry: same key, different body — should return error.
	tamperedBody := `{"amount":9999,"currency":"usd","payment_method":"pm_card_visa"}`
	_, _, err = store.Lookup(idempotencyKey, tamperedBody)
	if errors.Is(err, ErrHashMismatch) {
		fmt.Println("Hash mismatch detected — rejecting tampered retry (422)")
	}

	// Expiry check demo (entry is fresh, so not expired).
	expired := IsExpired(entry, 24)
	fmt.Printf("Entry expired (24h window): %v\n", expired)

	result := map[string]interface{}{
		"scenario":      "idempotent_retry_simulation",
		"first_stored":  stored,
		"cache_hit":     hit,
		"hash_mismatch": errors.Is(err, ErrHashMismatch),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
