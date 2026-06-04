package main

import "testing"

func TestValidateTilePolicyHealthy(t *testing.T) {
	cfg := TilePolicy{
		Name:               "healthy",
		ImmutableVersions:  true,
		ManifestTTLSeconds: 60,
		MaxActiveVersions:  4,
		PrewarmHotZooms:    true,
		RegionalCanary:     true,
		CDNMaxAgeSeconds:   86400,
	}
	if issues := ValidateTilePolicy(cfg); len(issues) != 0 {
		t.Fatalf("ValidateTilePolicy returned issues: %v", issues)
	}
}

func TestValidateTilePolicyWeak(t *testing.T) {
	cfg := TilePolicy{
		Name:               "weak",
		ManifestTTLSeconds: 0,
		MaxActiveVersions:  1,
		CDNMaxAgeSeconds:   60,
	}
	if issues := ValidateTilePolicy(cfg); len(issues) < 5 {
		t.Fatalf("ValidateTilePolicy returned too few issues: %v", issues)
	}
}
