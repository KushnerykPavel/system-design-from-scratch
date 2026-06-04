package main

import "testing"

func TestValidateLocationPipelineConfigHealthy(t *testing.T) {
	cfg := LocationPipelineConfig{
		Name:                  "healthy",
		Partitions:            256,
		DedupeWindowSeconds:   30,
		MaxAcceptedAgeSeconds: 15,
		MaxLateArrivalSeconds: 30,
		RegionalReplicas:      3,
		SmoothingEnabled:      true,
		LiveStateProjection:   true,
	}
	if issues := ValidateLocationPipelineConfig(cfg); len(issues) != 0 {
		t.Fatalf("ValidateLocationPipelineConfig returned issues: %v", issues)
	}
}

func TestValidateLocationPipelineConfigWeak(t *testing.T) {
	cfg := LocationPipelineConfig{
		Name:                  "weak",
		Partitions:            4,
		DedupeWindowSeconds:   0,
		MaxAcceptedAgeSeconds: 90,
		MaxLateArrivalSeconds: 10,
		RegionalReplicas:      1,
	}
	if issues := ValidateLocationPipelineConfig(cfg); len(issues) < 5 {
		t.Fatalf("ValidateLocationPipelineConfig returned too few issues: %v", issues)
	}
}
