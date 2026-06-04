package main

import "testing"

func TestValidateLogPipelineHealthy(t *testing.T) {
	p := LogPipeline{
		Name:                 "healthy",
		HotRetentionDays:     3,
		ArchiveRetentionDays: 30,
		PIIRedactionEnabled:  true,
		ReplayEnabled:        true,
		CriticalClassBuffer:  1,
		TenantQuotaEnabled:   true,
	}
	if issues := ValidateLogPipeline(p); len(issues) != 0 {
		t.Fatalf("ValidateLogPipeline returned issues: %v", issues)
	}
}

func TestValidateLogPipelineWeak(t *testing.T) {
	p := LogPipeline{
		Name:                 "weak",
		HotRetentionDays:     0,
		ArchiveRetentionDays: 0,
	}
	if issues := ValidateLogPipeline(p); len(issues) < 5 {
		t.Fatalf("ValidateLogPipeline returned too few issues: %v", issues)
	}
}
