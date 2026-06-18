package main

import (
	"testing"
	"time"
)

func TestWindowAggAcceptsInWindowEvents(t *testing.T) {
	now := time.Now()
	agg := NewWindowAgg(now.Add(-1*time.Minute), now)
	e := Event{UserID: "u1", TitleID: "tt-1", Type: "play", Timestamp: now.Add(-30 * time.Second)}
	accepted, _ := agg.Add(e, 10*time.Second)
	if !accepted {
		t.Fatal("expected event within window to be accepted")
	}
	if agg.Counts["play"] != 1 {
		t.Fatalf("expected play count 1, got %d", agg.Counts["play"])
	}
}

func TestWindowAggRejectsTooLateEvents(t *testing.T) {
	now := time.Now()
	agg := NewWindowAgg(now.Add(-1*time.Minute), now)
	// Event from 2 minutes ago is before the window start.
	e := Event{UserID: "u1", TitleID: "tt-1", Type: "play", Timestamp: now.Add(-2 * time.Minute)}
	accepted, late := agg.Add(e, 10*time.Second)
	if accepted {
		t.Fatal("expected event before window start to be rejected")
	}
	if !late {
		t.Fatal("expected rejection reason to be 'late'")
	}
}

func TestWindowAggAccumulatesTitleSignals(t *testing.T) {
	now := time.Now()
	agg := NewWindowAgg(now.Add(-1*time.Minute), now)
	e1 := Event{UserID: "u1", TitleID: "tt-1", Type: "end", Timestamp: now.Add(-30 * time.Second), Value: 90.0}
	e2 := Event{UserID: "u2", TitleID: "tt-1", Type: "end", Timestamp: now.Add(-20 * time.Second), Value: 50.0}
	agg.Add(e1, 10*time.Second)
	agg.Add(e2, 10*time.Second)
	if agg.Titles["tt-1"] != 140.0 {
		t.Fatalf("expected title signal 140.0, got %f", agg.Titles["tt-1"])
	}
}
