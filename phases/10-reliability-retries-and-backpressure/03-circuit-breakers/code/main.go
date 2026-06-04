package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BreakerConfig struct {
	TripErrorRate        float64 `json:"trip_error_rate"`
	TripTimeoutRate      float64 `json:"trip_timeout_rate"`
	UsesSaturationSignal bool    `json:"uses_saturation_signal"`
	FallbackIndependent  bool    `json:"fallback_independent"`
	HalfOpenMaxProbes    int     `json:"half_open_max_probes"`
	ScopedPerDependency  bool    `json:"scoped_per_dependency"`
}

type BreakerAssessment struct {
	Risk  string   `json:"risk"`
	Notes []string `json:"notes"`
}

func AssessBreaker(cfg BreakerConfig) BreakerAssessment {
	score := 0
	notes := make([]string, 0, 4)

	if !cfg.UsesSaturationSignal {
		score++
		notes = append(notes, "trip policy ignores saturation and may react too late")
	}
	if !cfg.FallbackIndependent {
		score += 2
		notes = append(notes, "fallback path depends on the same failure domain")
	}
	if cfg.HalfOpenMaxProbes > 10 || cfg.HalfOpenMaxProbes == 0 {
		score++
		notes = append(notes, "half-open probe volume is not tightly bounded")
	}
	if !cfg.ScopedPerDependency {
		score++
		notes = append(notes, "breaker scope is broad and may trip healthy traffic")
	}
	if cfg.TripErrorRate > 0.8 && cfg.TripTimeoutRate > 0.8 {
		score++
		notes = append(notes, "trip thresholds may be too tolerant for fast-failing protection")
	}

	risk := "low"
	if score >= 4 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return BreakerAssessment{Risk: risk, Notes: notes}
}

func main() {
	cfg := BreakerConfig{
		TripErrorRate:        0.35,
		TripTimeoutRate:      0.25,
		UsesSaturationSignal: true,
		FallbackIndependent:  true,
		HalfOpenMaxProbes:    5,
		ScopedPerDependency:  true,
	}

	encoded, err := json.MarshalIndent(AssessBreaker(cfg), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
