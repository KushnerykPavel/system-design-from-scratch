package main

import (
	"encoding/json"
	"os"
)

type ThreatModel struct {
	HasAssets        bool `json:"has_assets"`
	HasActors        bool `json:"has_actors"`
	HasBoundaries    bool `json:"has_boundaries"`
	HasTopThreats    bool `json:"has_top_threats"`
	HasMitigations   bool `json:"has_mitigations"`
	HasObservability bool `json:"has_observability"`
	HasDesignChange  bool `json:"has_design_change"`
}

type ThreatAssessment struct {
	Score   int      `json:"score"`
	Level   string   `json:"level"`
	Missing []string `json:"missing"`
}

func AssessThreatModel(model ThreatModel) ThreatAssessment {
	score := 0
	var missing []string

	add := func(ok bool, label string) {
		if ok {
			score += 2
			return
		}
		missing = append(missing, label)
	}

	add(model.HasAssets, "key assets")
	add(model.HasActors, "threat actors")
	add(model.HasBoundaries, "trust boundaries")
	add(model.HasTopThreats, "top prioritized threats")
	add(model.HasMitigations, "concrete mitigations")
	add(model.HasObservability, "security observability")
	add(model.HasDesignChange, "architecture change caused by threat model")

	level := "weak"
	switch {
	case score >= 12:
		level = "strong"
	case score >= 8:
		level = "developing"
	}

	return ThreatAssessment{Score: score, Level: level, Missing: missing}
}

func main() {
	model := ThreatModel{
		HasAssets:        true,
		HasActors:        true,
		HasBoundaries:    true,
		HasTopThreats:    true,
		HasMitigations:   true,
		HasObservability: false,
		HasDesignChange:  true,
	}
	_ = json.NewEncoder(os.Stdout).Encode(AssessThreatModel(model))
}
