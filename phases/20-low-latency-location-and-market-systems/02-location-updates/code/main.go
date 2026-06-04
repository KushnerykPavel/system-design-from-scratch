package main

import (
	"encoding/json"
	"flag"
	"os"
)

type LocationPipelineConfig struct {
	Name                  string `json:"name"`
	Partitions            int    `json:"partitions"`
	DedupeWindowSeconds   int    `json:"dedupe_window_seconds"`
	MaxAcceptedAgeSeconds int    `json:"max_accepted_age_seconds"`
	MaxLateArrivalSeconds int    `json:"max_late_arrival_seconds"`
	RegionalReplicas      int    `json:"regional_replicas"`
	SmoothingEnabled      bool   `json:"smoothing_enabled"`
	LiveStateProjection   bool   `json:"live_state_projection"`
}

func ValidateLocationPipelineConfig(cfg LocationPipelineConfig) []string {
	var issues []string
	if cfg.Partitions < 32 {
		issues = append(issues, "partitions should be at least 32 for a meaningful real-time pipeline")
	}
	if cfg.DedupeWindowSeconds <= 0 || cfg.DedupeWindowSeconds > 300 {
		issues = append(issues, "dedupe_window_seconds should stay between 1 and 300")
	}
	if cfg.MaxAcceptedAgeSeconds <= 0 || cfg.MaxAcceptedAgeSeconds > 60 {
		issues = append(issues, "max_accepted_age_seconds should reflect a tight live-serving freshness budget")
	}
	if cfg.MaxLateArrivalSeconds < cfg.MaxAcceptedAgeSeconds {
		issues = append(issues, "max_late_arrival_seconds should be at least max_accepted_age_seconds")
	}
	if cfg.RegionalReplicas < 2 {
		issues = append(issues, "regional_replicas should be at least 2")
	}
	if !cfg.SmoothingEnabled {
		issues = append(issues, "smoothing_enabled should usually be true for noisy GPS paths")
	}
	if !cfg.LiveStateProjection {
		issues = append(issues, "live_state_projection should be enabled for serving systems")
	}
	return issues
}

func main() {
	name := flag.String("name", "location-pipeline", "config name")
	flag.Parse()

	cfg := LocationPipelineConfig{
		Name:                  *name,
		Partitions:            256,
		DedupeWindowSeconds:   30,
		MaxAcceptedAgeSeconds: 15,
		MaxLateArrivalSeconds: 30,
		RegionalReplicas:      3,
		SmoothingEnabled:      true,
		LiveStateProjection:   true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"config": cfg,
		"issues": ValidateLocationPipelineConfig(cfg),
	})
}
