package main

import "testing"

func TestValidateRotationPlan(t *testing.T) {
	good := RotationPlan{
		ShortLivedCredentials: true,
		BackgroundRefresh:     true,
		OverlapWindow:         true,
		RevocationPath:        true,
	}
	if issues := ValidateRotationPlan(good); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	bad := RotationPlan{
		HotPathSecretFetch: true,
	}
	if issues := ValidateRotationPlan(bad); len(issues) < 4 {
		t.Fatalf("expected several issues, got %v", issues)
	}
}
