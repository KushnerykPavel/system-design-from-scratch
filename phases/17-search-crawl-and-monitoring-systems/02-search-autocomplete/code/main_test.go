package main

import "testing"

func TestValidateAutocompleteConfigHealthy(t *testing.T) {
	cfg := AutocompleteConfig{
		Name:               "healthy",
		MaxSuggestions:     8,
		FreshnessSeconds:   300,
		IndexReplicas:      3,
		PrefixCacheEnabled: true,
		PolicyFiltering:    true,
		FallbackSnapshot:   true,
	}
	if issues := ValidateAutocompleteConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateAutocompleteConfig returned issues: %v", issues)
	}
}

func TestValidateAutocompleteConfigWeak(t *testing.T) {
	cfg := AutocompleteConfig{
		Name:             "weak",
		MaxSuggestions:   30,
		FreshnessSeconds: 3600,
		IndexReplicas:    1,
	}
	if issues := ValidateAutocompleteConfig(cfg); len(issues) < 4 {
		t.Fatalf("ValidateAutocompleteConfig returned too few issues: %v", issues)
	}
}
