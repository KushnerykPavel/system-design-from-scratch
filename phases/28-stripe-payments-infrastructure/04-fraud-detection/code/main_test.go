package main

import (
	"math"
	"testing"
)

func TestLowRiskAllowed(t *testing.T) {
	tx := Transaction{
		ID: "txn_low", CardLast4: "4242",
		IPCountry: "US", BillingCountry: "US",
		AmountCents: 1000, IsNewCard: false,
	}
	signals := FraudSignals{
		VelocityScore: 0.05,
		GeoMismatch:   0.0,
		BINRisk:       0.05,
		DeviceRisk:    0.05,
	}
	score := ComputeRiskScore(tx, signals)
	action := RiskAction(score)
	if action != "allow" {
		t.Fatalf("expected allow for low-risk transaction, got %s (score=%.3f)", action, score)
	}
	if score >= 0.3 {
		t.Fatalf("expected score < 0.3 for low-risk transaction, got %.3f", score)
	}
}

func TestMediumRisk3DS(t *testing.T) {
	tx := Transaction{
		ID: "txn_med", CardLast4: "1234",
		IPCountry: "DE", BillingCountry: "US",
		AmountCents: 5000, IsNewCard: true,
	}
	signals := FraudSignals{
		VelocityScore: 0.20,
		GeoMismatch:   0.70,
		BINRisk:       0.20,
		DeviceRisk:    0.20,
	}
	score := ComputeRiskScore(tx, signals)
	action := RiskAction(score)
	if action != "3ds" {
		t.Fatalf("expected 3ds for medium-risk transaction, got %s (score=%.3f)", action, score)
	}
	if score < 0.3 || score >= 0.7 {
		t.Fatalf("expected 0.3 <= score < 0.7, got %.3f", score)
	}
}

func TestHighRiskBlocked(t *testing.T) {
	tx := Transaction{
		ID: "txn_high", CardLast4: "9999",
		IPCountry: "US", BillingCountry: "US",
		AmountCents: 50000, IsNewCard: false,
	}
	signals := FraudSignals{
		VelocityScore: 0.95,
		GeoMismatch:   0.10,
		BINRisk:       0.90,
		DeviceRisk:    0.90,
	}
	score := ComputeRiskScore(tx, signals)
	action := RiskAction(score)
	if action != "block" {
		t.Fatalf("expected block for high-risk transaction, got %s (score=%.3f)", action, score)
	}
	if score < 0.7 {
		t.Fatalf("expected score >= 0.7 for high-risk, got %.3f", score)
	}
}

func TestGeoMismatchImpact(t *testing.T) {
	base := FraudSignals{
		VelocityScore: 0.10,
		GeoMismatch:   0.0,
		BINRisk:       0.10,
		DeviceRisk:    0.10,
	}
	high := FraudSignals{
		VelocityScore: 0.10,
		GeoMismatch:   1.0,
		BINRisk:       0.10,
		DeviceRisk:    0.10,
	}
	tx := Transaction{ID: "txn_geo", AmountCents: 1000}
	scoreLow := ComputeRiskScore(tx, base)
	scoreHigh := ComputeRiskScore(tx, high)

	if scoreHigh <= scoreLow {
		t.Fatalf("geo mismatch should increase risk score: base=%.3f, high_geo=%.3f", scoreLow, scoreHigh)
	}
	// Difference should be 0.25 (geo weight = 0.25, delta geo = 1.0)
	expected := 0.25
	if math.Abs((scoreHigh-scoreLow)-expected) > 0.001 {
		t.Fatalf("expected geo mismatch delta of %.3f, got %.3f", expected, scoreHigh-scoreLow)
	}
}

func TestNewCardBump(t *testing.T) {
	signals := FraudSignals{
		VelocityScore: 0.10,
		GeoMismatch:   0.0,
		BINRisk:       0.10,
		DeviceRisk:    0.10,
	}
	existing := Transaction{ID: "txn_existing", AmountCents: 1000, IsNewCard: false}
	newCard := Transaction{ID: "txn_new", AmountCents: 1000, IsNewCard: true}

	scoreExisting := ComputeRiskScore(existing, signals)
	scoreNew := ComputeRiskScore(newCard, signals)

	if scoreNew <= scoreExisting {
		t.Fatalf("new card should have higher risk score: existing=%.3f, new=%.3f", scoreExisting, scoreNew)
	}
	diff := scoreNew - scoreExisting
	if math.Abs(diff-0.05) > 0.001 {
		t.Fatalf("expected new card bump of 0.05, got %.3f", diff)
	}
}

func TestScoreCapAt1(t *testing.T) {
	tx := Transaction{ID: "txn_max", AmountCents: 99999, IsNewCard: true}
	signals := FraudSignals{
		VelocityScore: 1.0,
		GeoMismatch:   1.0,
		BINRisk:       1.0,
		DeviceRisk:    1.0,
	}
	score := ComputeRiskScore(tx, signals)
	if score > 1.0 {
		t.Fatalf("risk score must not exceed 1.0, got %.3f", score)
	}
	if score != 1.0 {
		t.Fatalf("expected score=1.0 for all-max signals, got %.3f", score)
	}
}

func TestRiskActionBoundary(t *testing.T) {
	if RiskAction(0.0) != "allow" {
		t.Fatal("score 0.0 should allow")
	}
	if RiskAction(0.299) != "allow" {
		t.Fatal("score 0.299 should allow")
	}
	if RiskAction(0.3) != "3ds" {
		t.Fatal("score 0.3 should trigger 3ds")
	}
	if RiskAction(0.699) != "3ds" {
		t.Fatal("score 0.699 should trigger 3ds")
	}
	if RiskAction(0.7) != "block" {
		t.Fatal("score 0.7 should block")
	}
	if RiskAction(1.0) != "block" {
		t.Fatal("score 1.0 should block")
	}
}
