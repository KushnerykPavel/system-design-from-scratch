package main

import (
	"encoding/json"
	"os"
)

// Transaction represents an incoming payment to be fraud-scored.
type Transaction struct {
	ID             string
	CardLast4      string
	IPCountry      string
	BillingCountry string
	AmountCents    int
	IsNewCard      bool
}

// FraudSignals holds the normalized risk scores (0.0–1.0) for each signal category.
type FraudSignals struct {
	VelocityScore float64 // high velocity of txns from same card/IP/device
	GeoMismatch   float64 // mismatch between billing country and IP country
	BINRisk       float64 // risk from card type (prepaid, commercial, etc.)
	DeviceRisk    float64 // device fingerprint risk (new device, seen in fraud)
}

// ComputeRiskScore combines fraud signals into a single weighted risk score (0.0–1.0).
// Weights: Velocity 35%, GeoMismatch 25%, BINRisk 20%, DeviceRisk 20%.
func ComputeRiskScore(tx Transaction, signals FraudSignals) float64 {
	score := signals.VelocityScore*0.35 +
		signals.GeoMismatch*0.25 +
		signals.BINRisk*0.20 +
		signals.DeviceRisk*0.20

	// New card adds a flat risk bump of 0.05 (capped at 1.0).
	if tx.IsNewCard {
		score += 0.05
	}

	if score > 1.0 {
		score = 1.0
	}
	return score
}

// RiskAction returns the recommended action based on the risk score.
// Thresholds:  score < 0.3 → allow, 0.3 ≤ score < 0.7 → 3ds, score ≥ 0.7 → block.
func RiskAction(score float64) string {
	switch {
	case score >= 0.7:
		return "block"
	case score >= 0.3:
		return "3ds"
	default:
		return "allow"
	}
}

// ScoreResult bundles the transaction, signals, computed score, and action.
type ScoreResult struct {
	TransactionID string  `json:"transaction_id"`
	RiskScore     float64 `json:"risk_score"`
	Action        string  `json:"action"`
}

func main() {
	scenarios := []struct {
		tx      Transaction
		signals FraudSignals
		label   string
	}{
		{
			label: "low-risk domestic transaction",
			tx: Transaction{
				ID: "txn_001", CardLast4: "4242",
				IPCountry: "US", BillingCountry: "US",
				AmountCents: 2999, IsNewCard: false,
			},
			signals: FraudSignals{
				VelocityScore: 0.05,
				GeoMismatch:   0.0,
				BINRisk:       0.1,
				DeviceRisk:    0.05,
			},
		},
		{
			label: "medium-risk: geo mismatch + new card",
			tx: Transaction{
				ID: "txn_002", CardLast4: "1234",
				IPCountry: "RU", BillingCountry: "US",
				AmountCents: 15000, IsNewCard: true,
			},
			signals: FraudSignals{
				VelocityScore: 0.15,
				GeoMismatch:   0.80,
				BINRisk:       0.20,
				DeviceRisk:    0.30,
			},
		},
		{
			label: "high-risk: high velocity + suspicious device",
			tx: Transaction{
				ID: "txn_003", CardLast4: "9999",
				IPCountry: "US", BillingCountry: "US",
				AmountCents: 50000, IsNewCard: false,
			},
			signals: FraudSignals{
				VelocityScore: 0.95,
				GeoMismatch:   0.10,
				BINRisk:       0.80,
				DeviceRisk:    0.90,
			},
		},
		{
			label: "geo mismatch only (should elevate score)",
			tx: Transaction{
				ID: "txn_004", CardLast4: "5555",
				IPCountry: "CN", BillingCountry: "US",
				AmountCents: 7500, IsNewCard: false,
			},
			signals: FraudSignals{
				VelocityScore: 0.10,
				GeoMismatch:   0.90,
				BINRisk:       0.15,
				DeviceRisk:    0.10,
			},
		},
	}

	type output struct {
		Label         string  `json:"label"`
		TransactionID string  `json:"transaction_id"`
		RiskScore     float64 `json:"risk_score"`
		Action        string  `json:"action"`
	}

	var results []output
	for _, s := range scenarios {
		score := ComputeRiskScore(s.tx, s.signals)
		action := RiskAction(score)
		results = append(results, output{
			Label:         s.label,
			TransactionID: s.tx.ID,
			RiskScore:     score,
			Action:        action,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(results)
}
