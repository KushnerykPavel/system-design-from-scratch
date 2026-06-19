package main

import (
	"errors"
	"testing"
)

func TestStoreNewKey(t *testing.T) {
	s := NewIdempotencyStore()
	stored := s.Store("key-1", `{"amount":100}`, 200, `{"id":"pi_1"}`)
	if !stored {
		t.Fatal("expected Store to return true for a new key")
	}
}

func TestStoreDuplicateKey(t *testing.T) {
	s := NewIdempotencyStore()
	s.Store("key-1", `{"amount":100}`, 200, `{"id":"pi_1"}`)
	stored := s.Store("key-1", `{"amount":100}`, 200, `{"id":"pi_1"}`)
	if stored {
		t.Fatal("expected Store to return false for a duplicate key")
	}
}

func TestLookupCacheHit(t *testing.T) {
	s := NewIdempotencyStore()
	body := `{"amount":500,"currency":"usd"}`
	s.Store("key-2", body, 200, `{"status":"succeeded"}`)

	entry, hit, err := s.Lookup("key-2", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hit {
		t.Fatal("expected a cache hit")
	}
	if entry.ResponseCode != 200 {
		t.Fatalf("expected response code 200, got %d", entry.ResponseCode)
	}
}

func TestLookupMissingKey(t *testing.T) {
	s := NewIdempotencyStore()
	_, hit, err := s.Lookup("nonexistent-key", `{"amount":100}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hit {
		t.Fatal("expected no cache hit for missing key")
	}
}

func TestLookupHashMismatch(t *testing.T) {
	s := NewIdempotencyStore()
	s.Store("key-3", `{"amount":100}`, 200, `{"status":"succeeded"}`)

	_, _, err := s.Lookup("key-3", `{"amount":9999}`)
	if !errors.Is(err, ErrHashMismatch) {
		t.Fatalf("expected ErrHashMismatch, got %v", err)
	}
}

func TestIsExpiredFresh(t *testing.T) {
	s := NewIdempotencyStore()
	body := `{"amount":100}`
	s.Store("key-4", body, 200, `{}`)
	entry, _, _ := s.Lookup("key-4", body)
	if IsExpired(entry, 24) {
		t.Fatal("expected fresh entry to not be expired")
	}
}
