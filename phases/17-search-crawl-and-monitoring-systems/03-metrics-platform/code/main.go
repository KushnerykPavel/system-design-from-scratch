package main

import (
	"encoding/json"
	"flag"
	"os"
)

type MetricsPlan struct {
	Name                     string `json:"name"`
	HotRetentionDays         int    `json:"hot_retention_days"`
	ColdRetentionDays        int    `json:"cold_retention_days"`
	CardinalityBudget        int    `json:"cardinality_budget"`
	IngestReplicas           int    `json:"ingest_replicas"`
	DownsamplingEnabled      bool   `json:"downsampling_enabled"`
	QueryIsolationEnabled    bool   `json:"query_isolation_enabled"`
	RemoteWriteBufferMinutes int    `json:"remote_write_buffer_minutes"`
}

func ValidateMetricsPlan(plan MetricsPlan) []string {
	var issues []string
	if plan.HotRetentionDays <= 0 {
		issues = append(issues, "hot_retention_days must be positive")
	}
	if plan.ColdRetentionDays < plan.HotRetentionDays {
		issues = append(issues, "cold_retention_days should be at least hot_retention_days")
	}
	if plan.CardinalityBudget <= 0 {
		issues = append(issues, "cardinality_budget must be positive")
	}
	if plan.IngestReplicas < 2 {
		issues = append(issues, "ingest_replicas should be at least 2 for operational durability")
	}
	if !plan.DownsamplingEnabled {
		issues = append(issues, "downsampling_enabled should usually be true for long retention")
	}
	if !plan.QueryIsolationEnabled {
		issues = append(issues, "query_isolation_enabled should protect alerts from exploratory queries")
	}
	if plan.RemoteWriteBufferMinutes < 1 {
		issues = append(issues, "remote_write_buffer_minutes should tolerate short network interruptions")
	}
	return issues
}

func main() {
	name := flag.String("name", "multi-tenant-metrics", "metrics plan name")
	flag.Parse()

	plan := MetricsPlan{
		Name:                     *name,
		HotRetentionDays:         7,
		ColdRetentionDays:        180,
		CardinalityBudget:        1000000,
		IngestReplicas:           3,
		DownsamplingEnabled:      true,
		QueryIsolationEnabled:    true,
		RemoteWriteBufferMinutes: 5,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateMetricsPlan(plan),
	})
}
