package main

import "testing"

func TestAssessSecureDefaults(t *testing.T) {
	strong := AssessSecureDefaults(SecureDefaultsAnswer{
		DenyByDefault:       true,
		ShortLivedSecrets:   true,
		DefaultQuotas:       true,
		TenantScopedStorage: true,
		DeletionPolicy:      true,
		ExpiringOverrides:   true,
		DegradedModePlan:    true,
	})
	if strong.Level != "strong" {
		t.Fatalf("expected strong result, got %+v", strong)
	}

	weak := AssessSecureDefaults(SecureDefaultsAnswer{})
	if weak.Level == "strong" || len(weak.Missing) == 0 {
		t.Fatalf("expected missing items, got %+v", weak)
	}
}
