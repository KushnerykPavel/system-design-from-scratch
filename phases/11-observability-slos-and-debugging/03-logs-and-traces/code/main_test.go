package main

import "testing"

func TestAssessCorrelationLowRisk(t *testing.T) {
	got := AssessCorrelation([]Hop{
		{Name: "a", HasTraceID: true, HasRequestID: true, StructuredLogging: true},
		{Name: "b", HasTraceID: true, HasRequestID: true, StructuredLogging: true},
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}

func TestAssessCorrelationHighRiskOnBrokenPropagation(t *testing.T) {
	got := AssessCorrelation([]Hop{
		{Name: "gateway", HasTraceID: true, HasRequestID: true, StructuredLogging: true},
		{Name: "queue-worker", HasTraceID: false, HasRequestID: false, StructuredLogging: false},
		{Name: "notifier", HasTraceID: false, HasRequestID: true, StructuredLogging: false},
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}
