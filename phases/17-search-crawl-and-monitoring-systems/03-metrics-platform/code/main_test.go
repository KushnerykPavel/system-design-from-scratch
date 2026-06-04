package main

import "testing"

func TestValidateMetricsPlanHealthy(t *testing.T) {
	plan := MetricsPlan{
		Name:                     "healthy",
		HotRetentionDays:         7,
		ColdRetentionDays:        30,
		CardinalityBudget:        1000,
		IngestReplicas:           2,
		DownsamplingEnabled:      true,
		QueryIsolationEnabled:    true,
		RemoteWriteBufferMinutes: 2,
	}
	if issues := ValidateMetricsPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateMetricsPlan returned issues: %v", issues)
	}
}

func TestValidateMetricsPlanWeak(t *testing.T) {
	plan := MetricsPlan{
		Name:              "weak",
		HotRetentionDays:  0,
		ColdRetentionDays: 0,
		IngestReplicas:    1,
	}
	if issues := ValidateMetricsPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateMetricsPlan returned too few issues: %v", issues)
	}
}
