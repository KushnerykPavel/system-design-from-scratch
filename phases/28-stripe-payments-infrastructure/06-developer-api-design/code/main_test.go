package main

import (
	"testing"
	"time"
)

func TestBackoffSchedule(t *testing.T) {
	expected := []time.Duration{
		5 * time.Second,
		30 * time.Second,
		2 * time.Minute,
		10 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		3 * time.Hour,
		12 * time.Hour,
		24 * time.Hour,
	}
	for i, want := range expected {
		got := NextRetryDelay(i + 1)
		if got != want {
			t.Errorf("attempt %d: expected %s, got %s", i+1, want, got)
		}
	}
}

func TestBackoffCapAt24h(t *testing.T) {
	// Attempt 100 should still return 24h cap, not panic.
	got := NextRetryDelay(100)
	if got != 24*time.Hour {
		t.Fatalf("expected 24h cap for attempt 100, got %s", got)
	}
}

func TestBackoffAttemptZero(t *testing.T) {
	// Attempt 0 or negative should return first delay.
	if NextRetryDelay(0) != 5*time.Second {
		t.Fatalf("expected 5s for attempt 0, got %s", NextRetryDelay(0))
	}
	if NextRetryDelay(-1) != 5*time.Second {
		t.Fatalf("expected 5s for attempt -1, got %s", NextRetryDelay(-1))
	}
}

func TestSimulateDeliverySuccessAfterFailures(t *testing.T) {
	history := SimulateDelivery("evt_001", "https://example.com", 3)

	// Should have 4 entries: 3 failures + 1 success.
	if len(history) != 4 {
		t.Fatalf("expected 4 delivery entries, got %d", len(history))
	}

	// First 3 should be failed.
	for i := 0; i < 3; i++ {
		if history[i].Status != "failed" {
			t.Errorf("entry %d: expected status 'failed', got %q", i, history[i].Status)
		}
	}

	// Last entry should be delivered.
	last := history[len(history)-1]
	if last.Status != "delivered" {
		t.Fatalf("expected final status 'delivered', got %q", last.Status)
	}
	if last.Attempts != 4 {
		t.Fatalf("expected 4 total attempts, got %d", last.Attempts)
	}
}

func TestSimulateDeliveryImmediateSuccess(t *testing.T) {
	history := SimulateDelivery("evt_002", "https://example.com", 0)
	if len(history) != 1 {
		t.Fatalf("expected 1 delivery entry, got %d", len(history))
	}
	if history[0].Status != "delivered" {
		t.Fatalf("expected 'delivered', got %q", history[0].Status)
	}
}

func TestShouldRetryDelivered(t *testing.T) {
	d := WebhookDelivery{
		EventID:     "evt_003",
		Endpoint:    "https://example.com",
		Attempts:    1,
		NextRetryAt: time.Now().Add(5 * time.Second),
		Status:      "delivered",
	}
	if ShouldRetry(d, 9, 72) {
		t.Fatal("should not retry a delivered event")
	}
}

func TestShouldRetryAbandonedWindow(t *testing.T) {
	// Create a delivery that started more than 72 hours ago.
	d := WebhookDelivery{
		EventID:     "evt_004",
		Endpoint:    "https://example.com",
		Attempts:    5,
		NextRetryAt: time.Now().Add(-73 * time.Hour),
		Status:      "failed",
	}
	if ShouldRetry(d, 9, 72) {
		t.Fatal("should not retry after retry window has expired")
	}
}

func TestShouldRetryMaxAttemptsExceeded(t *testing.T) {
	d := WebhookDelivery{
		EventID:     "evt_005",
		Endpoint:    "https://example.com",
		Attempts:    9,
		NextRetryAt: time.Now().Add(1 * time.Hour),
		Status:      "failed",
	}
	if ShouldRetry(d, 9, 72) {
		t.Fatal("should not retry when maxAttempts reached")
	}
}
