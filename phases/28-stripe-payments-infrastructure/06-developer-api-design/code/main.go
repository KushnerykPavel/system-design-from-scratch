package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WebhookDelivery tracks the state of a single webhook delivery attempt.
type WebhookDelivery struct {
	EventID     string
	Endpoint    string
	Attempts    int
	NextRetryAt time.Time
	Status      string // "pending", "delivered", "failed", "abandoned"
}

// NextRetryDelay returns the exponential backoff delay for the given attempt number (1-indexed).
// Schedule: 5s, 30s, 2m, 10m, 30m, 1h, 3h, 12h, capped at 24h.
func NextRetryDelay(attempts int) time.Duration {
	delays := []time.Duration{
		5 * time.Second,
		30 * time.Second,
		2 * time.Minute,
		10 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		3 * time.Hour,
		12 * time.Hour,
		24 * time.Hour, // cap
	}
	if attempts <= 0 {
		return delays[0]
	}
	idx := attempts - 1
	if idx >= len(delays) {
		return delays[len(delays)-1]
	}
	return delays[idx]
}

// ShouldRetry returns true if the delivery should be retried.
// maxAttempts is the maximum number of delivery attempts.
// windowHours is the maximum retry window in hours (e.g., 72).
func ShouldRetry(d WebhookDelivery, maxAttempts int, windowHours int) bool {
	if d.Status == "delivered" || d.Status == "abandoned" {
		return false
	}
	if d.Attempts >= maxAttempts {
		return false
	}
	deadline := d.NextRetryAt.Add(-NextRetryDelay(d.Attempts)).Add(time.Duration(windowHours) * time.Hour)
	if time.Now().After(deadline) {
		return false
	}
	return true
}

// SimulateDelivery simulates a webhook delivery that fails the first `failFirst` attempts,
// then succeeds. Returns the sequence of WebhookDelivery states.
func SimulateDelivery(eventID, endpoint string, failFirst int) []WebhookDelivery {
	const maxAttempts = 9
	const windowHours = 72

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	var history []WebhookDelivery
	current := WebhookDelivery{
		EventID:     eventID,
		Endpoint:    endpoint,
		Attempts:    0,
		NextRetryAt: now,
		Status:      "pending",
	}

	for current.Attempts < maxAttempts {
		// Attempt delivery at NextRetryAt.
		current.Attempts++
		attemptTime := current.NextRetryAt

		if current.Attempts > failFirst {
			// Success.
			current.Status = "delivered"
			history = append(history, WebhookDelivery{
				EventID:     current.EventID,
				Endpoint:    current.Endpoint,
				Attempts:    current.Attempts,
				NextRetryAt: attemptTime,
				Status:      "delivered",
			})
			return history
		}

		// Failure: schedule next retry.
		nextDelay := NextRetryDelay(current.Attempts)
		nextRetry := attemptTime.Add(nextDelay)
		history = append(history, WebhookDelivery{
			EventID:     current.EventID,
			Endpoint:    current.Endpoint,
			Attempts:    current.Attempts,
			NextRetryAt: nextRetry,
			Status:      "failed",
		})

		// Check if we've exhausted the retry window.
		elapsed := nextRetry.Sub(now)
		if elapsed > time.Duration(windowHours)*time.Hour || current.Attempts >= maxAttempts {
			current.Status = "abandoned"
			history = append(history, WebhookDelivery{
				EventID:     current.EventID,
				Endpoint:    current.Endpoint,
				Attempts:    current.Attempts,
				NextRetryAt: nextRetry,
				Status:      "abandoned",
			})
			return history
		}
		current.NextRetryAt = nextRetry
	}

	current.Status = "abandoned"
	history = append(history, current)
	return history
}

func main() {
	fmt.Println("=== Webhook Retry Schedule (fails 3 times, then succeeds) ===")
	deliveries := SimulateDelivery("evt_test_001", "https://example.com/webhook", 3)

	type displayEntry struct {
		Attempt     int    `json:"attempt"`
		AttemptedAt string `json:"attempted_at"`
		Status      string `json:"status"`
		NextRetryIn string `json:"next_retry_in,omitempty"`
	}

	var display []displayEntry
	for i, d := range deliveries {
		entry := displayEntry{
			Attempt:     d.Attempts,
			AttemptedAt: d.NextRetryAt.Format(time.RFC3339),
			Status:      d.Status,
		}
		if d.Status == "failed" && i < len(deliveries)-1 {
			delay := NextRetryDelay(d.Attempts)
			entry.NextRetryIn = delay.String()
		}
		display = append(display, entry)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(display)

	fmt.Println("\n=== Backoff Schedule Preview ===")
	for i := 1; i <= 9; i++ {
		fmt.Printf("  attempt %d: retry after %s\n", i, NextRetryDelay(i))
	}
}
