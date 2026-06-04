package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Hop struct {
	Name              string `json:"name"`
	HasTraceID        bool   `json:"has_trace_id"`
	HasRequestID      bool   `json:"has_request_id"`
	StructuredLogging bool   `json:"structured_logging"`
}

type TraceAssessment struct {
	CoverageRatio float64  `json:"coverage_ratio"`
	Risk          string   `json:"risk"`
	Warnings      []string `json:"warnings"`
}

func AssessCorrelation(hops []Hop) TraceAssessment {
	if len(hops) == 0 {
		return TraceAssessment{Risk: "high", Warnings: []string{"no hops provided"}}
	}

	complete := 0
	warnings := make([]string, 0, len(hops))
	for _, hop := range hops {
		if hop.HasTraceID && hop.HasRequestID && hop.StructuredLogging {
			complete++
			continue
		}
		if !hop.HasTraceID {
			warnings = append(warnings, hop.Name+": missing trace propagation")
		}
		if !hop.HasRequestID {
			warnings = append(warnings, hop.Name+": missing request identifier")
		}
		if !hop.StructuredLogging {
			warnings = append(warnings, hop.Name+": unstructured logging")
		}
	}

	coverage := float64(complete) / float64(len(hops))
	risk := "low"
	if coverage < 0.7 {
		risk = "high"
	} else if coverage < 1 {
		risk = "medium"
	}

	return TraceAssessment{
		CoverageRatio: coverage,
		Risk:          risk,
		Warnings:      warnings,
	}
}

func main() {
	hops := []Hop{
		{Name: "gateway", HasTraceID: true, HasRequestID: true, StructuredLogging: true},
		{Name: "auth", HasTraceID: true, HasRequestID: true, StructuredLogging: true},
		{Name: "worker", HasTraceID: false, HasRequestID: true, StructuredLogging: false},
	}

	encoded, err := json.MarshalIndent(AssessCorrelation(hops), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
