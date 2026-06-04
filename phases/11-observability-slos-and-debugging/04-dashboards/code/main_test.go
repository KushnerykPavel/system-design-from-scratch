package main

import "testing"

func TestAssessDashboardPlanFlagsUnboundedLabels(t *testing.T) {
	got := AssessDashboardPlan(MetricPlan{
		BaseSeries: 10,
		Dimensions: []MetricDimension{
			{Name: "region", Cardinality: 20, Bounded: true},
			{Name: "tenant_id", Cardinality: 50000, Bounded: false},
		},
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessDashboardPlanKeepsBoundedTaxonomySafe(t *testing.T) {
	got := AssessDashboardPlan(MetricPlan{
		BaseSeries: 8,
		Dimensions: []MetricDimension{
			{Name: "region", Cardinality: 20, Bounded: true},
			{Name: "route_class", Cardinality: 10, Bounded: true},
			{Name: "status_class", Cardinality: 5, Bounded: true},
		},
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
