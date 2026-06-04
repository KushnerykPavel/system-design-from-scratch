package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ReservationPlan struct {
	Name                     string `json:"name"`
	AuthoritativeReservePath bool   `json:"authoritative_reserve_path"`
	ReservationTTL           bool   `json:"reservation_ttl"`
	IdempotentConfirm        bool   `json:"idempotent_confirm"`
	IdempotentRelease        bool   `json:"idempotent_release"`
	OversellGuard            bool   `json:"oversell_guard"`
	HotSKUProtection         bool   `json:"hot_sku_protection"`
	LeakDetection            bool   `json:"leak_detection"`
}

func ValidateReservationPlan(plan ReservationPlan) []string {
	var issues []string
	if !plan.AuthoritativeReservePath {
		issues = append(issues, "authoritative_reserve_path should be true so cached availability does not authorize reservations")
	}
	if !plan.ReservationTTL {
		issues = append(issues, "reservation_ttl should be true so abandoned holds are eventually reclaimed")
	}
	if !plan.IdempotentConfirm {
		issues = append(issues, "idempotent_confirm should be true so payment retries do not double-allocate stock")
	}
	if !plan.IdempotentRelease {
		issues = append(issues, "idempotent_release should be true so retries do not corrupt stock counters")
	}
	if !plan.OversellGuard {
		issues = append(issues, "oversell_guard should be true so concurrent reservations cannot exceed sellable inventory")
	}
	if !plan.HotSKUProtection {
		issues = append(issues, "hot_sku_protection should be true so flash traffic on one item does not collapse the system")
	}
	if !plan.LeakDetection {
		issues = append(issues, "leak_detection should be true so stuck reservations are surfaced quickly")
	}
	return issues
}

func main() {
	name := flag.String("name", "inventory-reservation", "plan name")
	flag.Parse()

	plan := ReservationPlan{
		Name:                     *name,
		AuthoritativeReservePath: true,
		ReservationTTL:           true,
		IdempotentConfirm:        true,
		IdempotentRelease:        true,
		OversellGuard:            true,
		HotSKUProtection:         true,
		LeakDetection:            true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateReservationPlan(plan),
	})
}
