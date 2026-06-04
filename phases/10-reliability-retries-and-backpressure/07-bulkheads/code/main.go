package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type IsolationPlan struct {
	DedicatedCriticalPool bool    `json:"dedicated_critical_pool"`
	TenantQuota           bool    `json:"tenant_quota"`
	CellCount             int     `json:"cell_count"`
	SharedCriticalDep     bool    `json:"shared_critical_dependency"`
	AllowsBorrowing       bool    `json:"allows_borrowing"`
	LargestTenantShare    float64 `json:"largest_tenant_share"`
}

type IsolationAssessment struct {
	Risk  string   `json:"risk"`
	Notes []string `json:"notes"`
}

func AssessIsolation(plan IsolationPlan) IsolationAssessment {
	score := 0
	notes := make([]string, 0, 4)

	if !plan.DedicatedCriticalPool {
		score++
		notes = append(notes, "critical work lacks dedicated execution capacity")
	}
	if plan.LargestTenantShare > 0.1 && !plan.TenantQuota {
		score += 2
		notes = append(notes, "large-tenant skew is unbounded")
	}
	if plan.CellCount < 2 {
		score++
		notes = append(notes, "there is no meaningful failure-domain partitioning")
	}
	if plan.SharedCriticalDep {
		score += 2
		notes = append(notes, "critical path still shares a single failure dependency")
	}
	if plan.AllowsBorrowing && !plan.DedicatedCriticalPool {
		score++
		notes = append(notes, "capacity borrowing exists without strong protected reservations")
	}

	risk := "low"
	if score >= 4 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return IsolationAssessment{Risk: risk, Notes: notes}
}

func main() {
	assessment := AssessIsolation(IsolationPlan{
		DedicatedCriticalPool: true,
		TenantQuota:           true,
		CellCount:             12,
		SharedCriticalDep:     false,
		AllowsBorrowing:       true,
		LargestTenantShare:    0.18,
	})

	encoded, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
