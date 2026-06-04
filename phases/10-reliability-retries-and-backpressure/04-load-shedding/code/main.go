package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type AdmissionPolicy struct {
	MaxInflight     int  `json:"max_inflight"`
	MaxQueueAgeMS   int  `json:"max_queue_age_ms"`
	HasPriorityLane bool `json:"has_priority_lane"`
	FastReject      bool `json:"fast_reject"`
	DownstreamAware bool `json:"downstream_aware"`
}

type Workload struct {
	PeakQPS       int     `json:"peak_qps"`
	SafeQPS       int     `json:"safe_qps"`
	CriticalShare float64 `json:"critical_share"`
}

type AdmissionAssessment struct {
	OverloadRatio float64  `json:"overload_ratio"`
	Risk          string   `json:"risk"`
	Notes         []string `json:"notes"`
}

func AssessAdmission(policy AdmissionPolicy, workload Workload) AdmissionAssessment {
	notes := make([]string, 0, 4)
	score := 0
	ratio := 0.0
	if workload.SafeQPS > 0 {
		ratio = float64(workload.PeakQPS) / float64(workload.SafeQPS)
	}
	if ratio > 1.2 && !policy.FastReject {
		score += 2
		notes = append(notes, "overload exceeds safe capacity without fast rejection")
	}
	if workload.CriticalShare > 0.1 && !policy.HasPriorityLane {
		score++
		notes = append(notes, "critical traffic lacks protected capacity")
	}
	if policy.MaxQueueAgeMS == 0 || policy.MaxQueueAgeMS > 5000 {
		score++
		notes = append(notes, "queue age limit is missing or too loose")
	}
	if !policy.DownstreamAware {
		score++
		notes = append(notes, "policy does not react to downstream saturation")
	}

	risk := "low"
	if score >= 4 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return AdmissionAssessment{OverloadRatio: ratio, Risk: risk, Notes: notes}
}

func main() {
	assessment := AssessAdmission(AdmissionPolicy{
		MaxInflight:     18000,
		MaxQueueAgeMS:   500,
		HasPriorityLane: true,
		FastReject:      true,
		DownstreamAware: true,
	}, Workload{
		PeakQPS:       250000,
		SafeQPS:       180000,
		CriticalShare: 0.2,
	})

	encoded, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
