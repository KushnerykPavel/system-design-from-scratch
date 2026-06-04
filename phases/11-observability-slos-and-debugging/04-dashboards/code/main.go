package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MetricDimension struct {
	Name        string `json:"name"`
	Cardinality int    `json:"cardinality"`
	Bounded     bool   `json:"bounded"`
}

type MetricPlan struct {
	BaseSeries int               `json:"base_series"`
	Dimensions []MetricDimension `json:"dimensions"`
}

type DashboardAssessment struct {
	EstimatedSeries int      `json:"estimated_series"`
	Risk            string   `json:"risk"`
	Warnings        []string `json:"warnings"`
}

func AssessDashboardPlan(plan MetricPlan) DashboardAssessment {
	estimate := max(plan.BaseSeries, 1)
	warnings := make([]string, 0, len(plan.Dimensions))
	unbounded := false

	for _, dim := range plan.Dimensions {
		cardinality := dim.Cardinality
		if cardinality < 1 {
			cardinality = 1
		}
		estimate *= cardinality
		if !dim.Bounded {
			unbounded = true
			warnings = append(warnings, dim.Name+": unbounded label")
		}
		if cardinality > 1000 {
			warnings = append(warnings, dim.Name+": very high label cardinality")
		}
	}

	risk := "low"
	switch {
	case unbounded || estimate > 1000000:
		risk = "high"
	case estimate > 100000:
		risk = "medium"
	}

	return DashboardAssessment{
		EstimatedSeries: estimate,
		Risk:            risk,
		Warnings:        warnings,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	plan := MetricPlan{
		BaseSeries: 12,
		Dimensions: []MetricDimension{
			{Name: "region", Cardinality: 20, Bounded: true},
			{Name: "route_class", Cardinality: 12, Bounded: true},
			{Name: "tenant_id", Cardinality: 50000, Bounded: false},
		},
	}

	encoded, err := json.MarshalIndent(AssessDashboardPlan(plan), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
