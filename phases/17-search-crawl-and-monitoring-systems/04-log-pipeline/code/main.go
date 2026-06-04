package main

import (
	"encoding/json"
	"flag"
	"os"
)

type LogPipeline struct {
	Name                 string `json:"name"`
	HotRetentionDays     int    `json:"hot_retention_days"`
	ArchiveRetentionDays int    `json:"archive_retention_days"`
	PIIRedactionEnabled  bool   `json:"pii_redaction_enabled"`
	ReplayEnabled        bool   `json:"replay_enabled"`
	CriticalClassBuffer  int    `json:"critical_class_buffer"`
	TenantQuotaEnabled   bool   `json:"tenant_quota_enabled"`
}

func ValidateLogPipeline(p LogPipeline) []string {
	var issues []string
	if p.HotRetentionDays <= 0 {
		issues = append(issues, "hot_retention_days must be positive")
	}
	if p.ArchiveRetentionDays < p.HotRetentionDays {
		issues = append(issues, "archive_retention_days should be at least hot_retention_days")
	}
	if !p.PIIRedactionEnabled {
		issues = append(issues, "pii_redaction_enabled should be true before indexing logs")
	}
	if !p.ReplayEnabled {
		issues = append(issues, "replay_enabled should usually be true for parser recovery")
	}
	if p.CriticalClassBuffer < 1 {
		issues = append(issues, "critical_class_buffer should reserve capacity for protected log classes")
	}
	if !p.TenantQuotaEnabled {
		issues = append(issues, "tenant_quota_enabled should be true in multi-tenant log platforms")
	}
	return issues
}

func main() {
	name := flag.String("name", "global-log-pipeline", "log pipeline name")
	flag.Parse()

	p := LogPipeline{
		Name:                 *name,
		HotRetentionDays:     3,
		ArchiveRetentionDays: 180,
		PIIRedactionEnabled:  true,
		ReplayEnabled:        true,
		CriticalClassBuffer:  2,
		TenantQuotaEnabled:   true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"pipeline": p,
		"issues":   ValidateLogPipeline(p),
	})
}
