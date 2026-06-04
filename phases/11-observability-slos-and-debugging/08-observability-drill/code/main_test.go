package main

import "testing"

func TestAssessDrillStrongAnswer(t *testing.T) {
	got := AssessDrill(DrillAnswer{
		HasSLI:        true,
		HasSLO:        true,
		HasMetrics:    true,
		HasLogsTraces: true,
		HasDashboards: true,
		HasAlerts:     true,
		HasRunbook:    true,
		HasDebugStory: true,
		HasTradeoffs:  true,
	})

	if got.Level != "strong" {
		t.Fatalf("level = %q, want strong", got.Level)
	}
}

func TestAssessDrillFlagsMissingCoreAreas(t *testing.T) {
	got := AssessDrill(DrillAnswer{
		HasSLI:       true,
		HasSLO:       false,
		HasMetrics:   true,
		HasTradeoffs: false,
	})

	if got.Level == "strong" {
		t.Fatalf("level = %q, want not strong", got.Level)
	}
	if len(got.Missing) == 0 {
		t.Fatalf("expected missing areas")
	}
}
