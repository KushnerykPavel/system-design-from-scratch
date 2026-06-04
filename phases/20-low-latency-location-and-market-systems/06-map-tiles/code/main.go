package main

import (
	"encoding/json"
	"flag"
	"os"
)

type TilePolicy struct {
	Name               string `json:"name"`
	ImmutableVersions  bool   `json:"immutable_versions"`
	ManifestTTLSeconds int    `json:"manifest_ttl_seconds"`
	MaxActiveVersions  int    `json:"max_active_versions"`
	PrewarmHotZooms    bool   `json:"prewarm_hot_zooms"`
	RegionalCanary     bool   `json:"regional_canary"`
	CDNMaxAgeSeconds   int    `json:"cdn_max_age_seconds"`
}

func ValidateTilePolicy(cfg TilePolicy) []string {
	var issues []string
	if !cfg.ImmutableVersions {
		issues = append(issues, "immutable_versions should be enabled")
	}
	if cfg.ManifestTTLSeconds <= 0 || cfg.ManifestTTLSeconds > 3600 {
		issues = append(issues, "manifest_ttl_seconds should stay between 1 and 3600")
	}
	if cfg.MaxActiveVersions < 2 || cfg.MaxActiveVersions > 20 {
		issues = append(issues, "max_active_versions should stay between 2 and 20")
	}
	if !cfg.PrewarmHotZooms {
		issues = append(issues, "prewarm_hot_zooms should usually be enabled")
	}
	if !cfg.RegionalCanary {
		issues = append(issues, "regional_canary should be enabled")
	}
	if cfg.CDNMaxAgeSeconds < 300 {
		issues = append(issues, "cdn_max_age_seconds should be at least 300 for read-mostly immutable assets")
	}
	return issues
}

func main() {
	name := flag.String("name", "map-tiles", "policy name")
	flag.Parse()

	cfg := TilePolicy{
		Name:               *name,
		ImmutableVersions:  true,
		ManifestTTLSeconds: 60,
		MaxActiveVersions:  4,
		PrewarmHotZooms:    true,
		RegionalCanary:     true,
		CDNMaxAgeSeconds:   86400,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateTilePolicy(cfg),
	})
}
