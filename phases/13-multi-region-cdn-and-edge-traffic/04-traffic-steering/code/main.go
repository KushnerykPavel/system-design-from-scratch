package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type SteeringPolicy struct {
	Name                   string `json:"name"`
	IngressMode            string `json:"ingress_mode"`
	PropagationSeconds     int    `json:"propagation_seconds"`
	HealthSignalWindowSecs int    `json:"health_signal_window_secs"`
	HasCapacitySignals     bool   `json:"has_capacity_signals"`
	SupportsAffinity       bool   `json:"supports_affinity"`
	HasEmergencyOverride   bool   `json:"has_emergency_override"`
	HasRouteExplanation    bool   `json:"has_route_explanation"`
}

func ValidatePolicy(policy SteeringPolicy) []string {
	var issues []string
	if policy.IngressMode != "dns" && policy.IngressMode != "anycast" {
		issues = append(issues, "ingress_mode must be dns or anycast")
	}
	if policy.PropagationSeconds <= 0 || policy.PropagationSeconds > 300 {
		issues = append(issues, "propagation_seconds should be between 1 and 300")
	}
	if policy.HealthSignalWindowSecs < 10 {
		issues = append(issues, "health signal window is too short and may cause flapping")
	}
	if !policy.HasCapacitySignals {
		issues = append(issues, "routing policy should consider downstream capacity signals")
	}
	if !policy.HasEmergencyOverride {
		issues = append(issues, "emergency override path is required")
	}
	if !policy.HasRouteExplanation {
		issues = append(issues, "route decisions should be explainable")
	}
	return issues
}

func main() {
	name := flag.String("name", "global-routing", "steering policy name")
	flag.Parse()

	policy := SteeringPolicy{
		Name:                   *name,
		IngressMode:            "anycast",
		PropagationSeconds:     30,
		HealthSignalWindowSecs: 20,
		HasCapacitySignals:     true,
		SupportsAffinity:       true,
		HasEmergencyOverride:   true,
		HasRouteExplanation:    true,
	}

	payload := map[string]any{
		"policy": policy,
		"issues": ValidatePolicy(policy),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
