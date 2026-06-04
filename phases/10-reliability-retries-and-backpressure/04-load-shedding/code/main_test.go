package main

import "testing"

func TestAssessAdmissionFlagsUnprotectedOverload(t *testing.T) {
	got := AssessAdmission(AdmissionPolicy{
		MaxInflight:     40000,
		MaxQueueAgeMS:   15000,
		HasPriorityLane: false,
		FastReject:      false,
		DownstreamAware: false,
	}, Workload{
		PeakQPS:       250000,
		SafeQPS:       150000,
		CriticalShare: 0.2,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessAdmissionApprovesBoundedPriorityPolicy(t *testing.T) {
	got := AssessAdmission(AdmissionPolicy{
		MaxInflight:     18000,
		MaxQueueAgeMS:   500,
		HasPriorityLane: true,
		FastReject:      true,
		DownstreamAware: true,
	}, Workload{
		PeakQPS:       180000,
		SafeQPS:       180000,
		CriticalShare: 0.25,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
