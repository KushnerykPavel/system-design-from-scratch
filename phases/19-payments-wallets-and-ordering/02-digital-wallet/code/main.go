package main

import (
	"encoding/json"
	"flag"
	"os"
)

type WalletPlan struct {
	Name                    string `json:"name"`
	TracksAvailableAndHeld  bool   `json:"tracks_available_and_held"`
	HasHoldExpiry           bool   `json:"has_hold_expiry"`
	IdempotentSettlement    bool   `json:"idempotent_settlement"`
	IdempotentRelease       bool   `json:"idempotent_release"`
	PreventsNegativeBalance bool   `json:"prevents_negative_balance"`
	SupportsPartialSettle   bool   `json:"supports_partial_settle"`
	HasAuditTrail           bool   `json:"has_audit_trail"`
}

func ValidateWalletPlan(plan WalletPlan) []string {
	var issues []string
	if !plan.TracksAvailableAndHeld {
		issues = append(issues, "tracks_available_and_held should be true so reserved funds stay distinct from spendable funds")
	}
	if !plan.HasHoldExpiry {
		issues = append(issues, "has_hold_expiry should be true so abandoned holds do not reduce balance forever")
	}
	if !plan.IdempotentSettlement {
		issues = append(issues, "idempotent_settlement should be true so order retries do not double-settle holds")
	}
	if !plan.IdempotentRelease {
		issues = append(issues, "idempotent_release should be true so cancellation retries stay safe")
	}
	if !plan.PreventsNegativeBalance {
		issues = append(issues, "prevents_negative_balance should be true so the wallet cannot overspend")
	}
	if !plan.SupportsPartialSettle {
		issues = append(issues, "supports_partial_settle should be true so final amount adjustments are expressible")
	}
	if !plan.HasAuditTrail {
		issues = append(issues, "has_audit_trail should be true so support and finance can inspect balance changes")
	}
	return issues
}

func main() {
	name := flag.String("name", "digital-wallet", "plan name")
	flag.Parse()

	plan := WalletPlan{
		Name:                    *name,
		TracksAvailableAndHeld:  true,
		HasHoldExpiry:           true,
		IdempotentSettlement:    true,
		IdempotentRelease:       true,
		PreventsNegativeBalance: true,
		SupportsPartialSettle:   true,
		HasAuditTrail:           true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateWalletPlan(plan),
	})
}
