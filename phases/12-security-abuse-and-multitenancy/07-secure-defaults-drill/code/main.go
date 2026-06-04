package main

import (
	"encoding/json"
	"os"
)

type SecureDefaultsAnswer struct {
	DenyByDefault       bool `json:"deny_by_default"`
	ShortLivedSecrets   bool `json:"short_lived_secrets"`
	DefaultQuotas       bool `json:"default_quotas"`
	TenantScopedStorage bool `json:"tenant_scoped_storage"`
	DeletionPolicy      bool `json:"deletion_policy"`
	ExpiringOverrides   bool `json:"expiring_overrides"`
	DegradedModePlan    bool `json:"degraded_mode_plan"`
}

type DrillResult struct {
	Score   int      `json:"score"`
	Level   string   `json:"level"`
	Missing []string `json:"missing"`
}

func AssessSecureDefaults(answer SecureDefaultsAnswer) DrillResult {
	score := 0
	var missing []string

	add := func(ok bool, label string) {
		if ok {
			score += 2
			return
		}
		missing = append(missing, label)
	}

	add(answer.DenyByDefault, "deny-by-default policy")
	add(answer.ShortLivedSecrets, "managed short-lived secrets")
	add(answer.DefaultQuotas, "default quotas or abuse controls")
	add(answer.TenantScopedStorage, "tenant-scoped storage")
	add(answer.DeletionPolicy, "deletion and retention policy")
	add(answer.ExpiringOverrides, "expiring override workflow")
	add(answer.DegradedModePlan, "degraded-mode security behavior")

	level := "weak"
	switch {
	case score >= 12:
		level = "strong"
	case score >= 8:
		level = "developing"
	}

	return DrillResult{Score: score, Level: level, Missing: missing}
}

func main() {
	answer := SecureDefaultsAnswer{
		DenyByDefault:       true,
		ShortLivedSecrets:   true,
		DefaultQuotas:       true,
		TenantScopedStorage: true,
		DeletionPolicy:      true,
		ExpiringOverrides:   false,
		DegradedModePlan:    true,
	}
	_ = json.NewEncoder(os.Stdout).Encode(AssessSecureDefaults(answer))
}
