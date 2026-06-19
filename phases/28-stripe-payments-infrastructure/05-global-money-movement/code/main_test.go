package main

import (
	"strings"
	"testing"
)

var testRates = map[string]float64{
	"USD/EUR": 0.92,
	"USD/JPY": 149.50,
	"EUR/USD": 1.087,
	"USD/BHD": 0.377,
}

func TestConvertUSDtoEUR(t *testing.T) {
	from := Money{AmountCents: 1000, Currency: "USD"} // $10.00
	result, err := ConvertCurrency(from, "EUR", testRates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Currency != "EUR" {
		t.Fatalf("expected EUR, got %s", result.Currency)
	}
	// $10.00 * 0.92 = €9.20 → 920 euro cents
	if result.AmountCents != 920 {
		t.Fatalf("expected 920 EUR cents, got %d", result.AmountCents)
	}
}

func TestConvertUSDtoJPY_ZeroDecimal(t *testing.T) {
	from := Money{AmountCents: 1000, Currency: "USD"} // $10.00
	result, err := ConvertCurrency(from, "JPY", testRates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Currency != "JPY" {
		t.Fatalf("expected JPY, got %s", result.Currency)
	}
	// $10.00 * 149.50 = ¥1495
	if result.AmountCents != 1495 {
		t.Fatalf("expected 1495 JPY, got %d", result.AmountCents)
	}
}

func TestConvertUSDtoBHD_ThreeDecimal(t *testing.T) {
	from := Money{AmountCents: 10000, Currency: "USD"} // $100.00
	result, err := ConvertCurrency(from, "BHD", testRates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Currency != "BHD" {
		t.Fatalf("expected BHD, got %s", result.Currency)
	}
	// $100.00 * 0.377 = 37.700 BHD → 37700 fils
	if result.AmountCents != 37700 {
		t.Fatalf("expected 37700 BHD fils, got %d", result.AmountCents)
	}
}

func TestConvertSameCurrency(t *testing.T) {
	from := Money{AmountCents: 5000, Currency: "JPY"}
	result, err := ConvertCurrency(from, "JPY", testRates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AmountCents != 5000 {
		t.Fatalf("expected 5000 JPY unchanged, got %d", result.AmountCents)
	}
}

func TestConvertMissingRate(t *testing.T) {
	from := Money{AmountCents: 1000, Currency: "USD"}
	_, err := ConvertCurrency(from, "CHF", testRates)
	if err == nil {
		t.Fatal("expected error for missing rate, got nil")
	}
}

func TestValidateMoneyNegative(t *testing.T) {
	err := ValidateMoney(Money{AmountCents: -1, Currency: "USD"})
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
}

func TestValidateMoneyInvalidCode(t *testing.T) {
	err := ValidateMoney(Money{AmountCents: 100, Currency: "US"})
	if err == nil {
		t.Fatal("expected error for 2-letter currency code")
	}
}

func TestValidateMoneyValid(t *testing.T) {
	if err := ValidateMoney(Money{AmountCents: 0, Currency: "USD"}); err != nil {
		t.Fatalf("unexpected error for zero amount: %v", err)
	}
	if err := ValidateMoney(Money{AmountCents: 1000, Currency: "JPY"}); err != nil {
		t.Fatalf("unexpected error for JPY: %v", err)
	}
}

func TestFormatMoneyUSD(t *testing.T) {
	result := FormatMoney(Money{AmountCents: 1099, Currency: "USD"})
	if result != "10.99 USD" {
		t.Fatalf("expected '10.99 USD', got %q", result)
	}
}

func TestFormatMoneyJPY_ZeroDecimal(t *testing.T) {
	result := FormatMoney(Money{AmountCents: 1000, Currency: "JPY"})
	if result != "1000 JPY" {
		t.Fatalf("expected '1000 JPY', got %q", result)
	}
	// JPY format must NOT contain a decimal point
	if strings.Contains(result, ".") {
		t.Fatalf("JPY format should not contain decimal point, got %q", result)
	}
}

func TestFormatMoneyBHD_ThreeDecimal(t *testing.T) {
	result := FormatMoney(Money{AmountCents: 10000, Currency: "BHD"})
	if result != "10.000 BHD" {
		t.Fatalf("expected '10.000 BHD', got %q", result)
	}
}

func TestFormatMoneyEURSmall(t *testing.T) {
	result := FormatMoney(Money{AmountCents: 5, Currency: "EUR"})
	if result != "0.05 EUR" {
		t.Fatalf("expected '0.05 EUR', got %q", result)
	}
}
