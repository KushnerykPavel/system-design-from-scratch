package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Event represents a single playback or UI event from a device.
type Event struct {
	UserID    string    `json:"user_id"`
	TitleID   string    `json:"title_id"`
	Type      string    `json:"type"` // "play", "heartbeat", "end"
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"` // completion_percent for "end", buffer_ms for "heartbeat"
}

// WindowAgg accumulates events within a sliding window.
type WindowAgg struct {
	WindowStart time.Time
	WindowEnd   time.Time
	Counts      map[string]int     // event type -> count
	Titles      map[string]float64 // title_id -> total completion
}

// NewWindowAgg creates an aggregation window.
func NewWindowAgg(start, end time.Time) *WindowAgg {
	return &WindowAgg{
		WindowStart: start,
		WindowEnd:   end,
		Counts:      make(map[string]int),
		Titles:      make(map[string]float64),
	}
}

// Add incorporates an event into the window if it falls within the time range.
// Returns true if the event was accepted, false if it was late.
func (w *WindowAgg) Add(e Event, watermarkLag time.Duration) (accepted bool, late bool) {
	watermark := w.WindowEnd.Add(-watermarkLag)
	if e.Timestamp.Before(w.WindowStart) {
		return false, true // too late
	}
	if e.Timestamp.After(w.WindowEnd) {
		return false, false // too early for this window
	}
	if e.Timestamp.Before(watermark) && e.Timestamp.After(w.WindowStart) {
		late = true
	}
	w.Counts[e.Type]++
	if e.Type == "end" {
		w.Titles[e.TitleID] += e.Value
	}
	return true, late
}

// Summary returns the aggregated result for this window.
func (w *WindowAgg) Summary() map[string]any {
	return map[string]any{
		"window_start":  w.WindowStart,
		"window_end":    w.WindowEnd,
		"event_counts":  w.Counts,
		"title_signals": w.Titles,
	}
}

func main() {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	windowStart := now.Add(-1 * time.Minute)
	windowEnd := now
	watermarkLag := 10 * time.Second

	agg := NewWindowAgg(windowStart, windowEnd)

	events := []Event{
		{UserID: "u1", TitleID: "tt-001", Type: "play", Timestamp: now.Add(-50 * time.Second)},
		{UserID: "u2", TitleID: "tt-002", Type: "end", Timestamp: now.Add(-30 * time.Second), Value: 95.0},
		{UserID: "u3", TitleID: "tt-001", Type: "end", Timestamp: now.Add(-5 * time.Second), Value: 12.0},
		// Late event: arrived within window time but past watermark would have closed it.
		{UserID: "u4", TitleID: "tt-003", Type: "end", Timestamp: now.Add(-55 * time.Second), Value: 88.0},
	}

	lateCount := 0
	for _, e := range events {
		accepted, late := agg.Add(e, watermarkLag)
		if accepted && late {
			lateCount++
		}
	}

	result := map[string]any{
		"window_summary": agg.Summary(),
		"late_events":    lateCount,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
