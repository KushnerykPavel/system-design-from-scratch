package main

import (
	"encoding/json"
	"flag"
	"os"
)

type RedirectPlan struct {
	Name                  string `json:"name"`
	CodeLength            int    `json:"code_length"`
	CacheTTLSeconds       int    `json:"cache_ttl_seconds"`
	UsesEdgeCache         bool   `json:"uses_edge_cache"`
	HasIdempotencyKeys    bool   `json:"has_idempotency_keys"`
	AsyncAnalytics        bool   `json:"async_analytics"`
	SupportsCustomAlias   bool   `json:"supports_custom_alias"`
	HotKeyMitigation      bool   `json:"hot_key_mitigation"`
	PolicyScanBeforeWrite bool   `json:"policy_scan_before_write"`
}

func ValidateRedirectPlan(plan RedirectPlan) []string {
	var issues []string
	if plan.CodeLength < 6 {
		issues = append(issues, "code_length should leave enough namespace headroom for sustained growth")
	}
	if plan.CacheTTLSeconds <= 0 {
		issues = append(issues, "cache_ttl_seconds must be positive for redirect-heavy workloads")
	}
	if !plan.UsesEdgeCache {
		issues = append(issues, "uses_edge_cache should be true when redirects dominate traffic")
	}
	if !plan.HasIdempotencyKeys {
		issues = append(issues, "has_idempotency_keys should be true to keep retries from creating duplicate links")
	}
	if !plan.AsyncAnalytics {
		issues = append(issues, "async_analytics should be true unless the product explicitly requires synchronous accounting")
	}
	if plan.SupportsCustomAlias && !plan.PolicyScanBeforeWrite {
		issues = append(issues, "policy_scan_before_write should be enabled before activating user-controlled aliases")
	}
	if !plan.HotKeyMitigation {
		issues = append(issues, "hot_key_mitigation should be enabled for viral link skew")
	}
	return issues
}

func main() {
	name := flag.String("name", "global-shortener", "plan name")
	flag.Parse()

	plan := RedirectPlan{
		Name:                  *name,
		CodeLength:            8,
		CacheTTLSeconds:       300,
		UsesEdgeCache:         true,
		HasIdempotencyKeys:    true,
		AsyncAnalytics:        true,
		SupportsCustomAlias:   true,
		HotKeyMitigation:      true,
		PolicyScanBeforeWrite: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateRedirectPlan(plan),
	})
}
