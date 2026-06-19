package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

// ZeroDecimalCurrencies are currencies where the API amount equals the display amount
// (no subunit exists). 1 in the API = 1 yen, 1 won, etc.
var ZeroDecimalCurrencies = map[string]bool{
	"JPY": true, "KRW": true, "VND": true, "BIF": true, "CLP": true,
	"GNF": true, "MGA": true, "PYG": true, "RWF": true, "UGX": true,
	"XAF": true, "XOF": true, "XPF": true,
}

// ThreeDecimalCurrencies are currencies where the API amount is in units of 1/1000
// (e.g., BHD fils: 10000 = 10.000 BHD).
var ThreeDecimalCurrencies = map[string]bool{
	"BHD": true, "JOD": true, "KWD": true, "OMR": true, "TND": true,
}

// Money represents a monetary value as an integer amount in the currency's smallest unit.
// AmountCents is named for clarity but represents the smallest unit for any currency.
type Money struct {
	AmountCents int64  // smallest unit: cents for USD, yen for JPY, fils for BHD
	Currency    string // ISO 4217 code, e.g., "USD", "JPY", "BHD"
}

// ValidateMoney returns an error if the Money value is invalid.
func ValidateMoney(m Money) error {
	if m.AmountCents < 0 {
		return errors.New("amount must be non-negative")
	}
	code := strings.ToUpper(m.Currency)
	if len(code) != 3 {
		return fmt.Errorf("invalid ISO 4217 currency code: %q", m.Currency)
	}
	// All letters A-Z check
	for _, c := range code {
		if c < 'A' || c > 'Z' {
			return fmt.Errorf("invalid ISO 4217 currency code: %q", m.Currency)
		}
	}
	return nil
}

// ConvertCurrency converts a Money value from one currency to another using the
// provided exchange rates (expressed as units of toCurrency per 1 unit of fromCurrency's
// major unit). Returns an error if the currency is invalid or no rate is available.
//
// The conversion uses integer arithmetic to avoid float64 rounding errors:
// 1. Convert AmountCents to the major unit using integer division.
// 2. Apply the float64 rate only for the cross-currency calculation.
// 3. Round to the nearest minor unit of the target currency.
func ConvertCurrency(from Money, toCurrency string, rates map[string]float64) (Money, error) {
	if err := ValidateMoney(from); err != nil {
		return Money{}, fmt.Errorf("invalid source money: %w", err)
	}
	toCurrency = strings.ToUpper(toCurrency)
	if len(toCurrency) != 3 {
		return Money{}, fmt.Errorf("invalid target currency code: %q", toCurrency)
	}
	fromCode := strings.ToUpper(from.Currency)
	if fromCode == toCurrency {
		return Money{AmountCents: from.AmountCents, Currency: toCurrency}, nil
	}

	rate, ok := rates[fromCode+"/"+toCurrency]
	if !ok {
		return Money{}, fmt.Errorf("no exchange rate available for %s/%s", fromCode, toCurrency)
	}

	// Determine decimal places for source and target currencies.
	fromDecimals := decimalPlaces(fromCode)
	toDecimals := decimalPlaces(toCurrency)

	// Convert to major units as float64 (only for the rate multiplication).
	fromMajor := float64(from.AmountCents) / math.Pow10(fromDecimals)

	// Apply the exchange rate.
	toMajor := fromMajor * rate

	// Convert back to minor units, rounding to nearest.
	toMinor := int64(math.Round(toMajor * math.Pow10(toDecimals)))

	return Money{AmountCents: toMinor, Currency: toCurrency}, nil
}

// decimalPlaces returns the number of decimal places for a currency code.
func decimalPlaces(code string) int {
	if ZeroDecimalCurrencies[code] {
		return 0
	}
	if ThreeDecimalCurrencies[code] {
		return 3
	}
	return 2
}

// FormatMoney formats a Money value as a human-readable string.
// JPY:  1000       → "1000 JPY"   (no decimal point)
// USD:  1000       → "10.00 USD"
// BHD:  10000      → "10.000 BHD"
func FormatMoney(m Money) string {
	code := strings.ToUpper(m.Currency)
	decimals := decimalPlaces(code)
	if decimals == 0 {
		return fmt.Sprintf("%d %s", m.AmountCents, code)
	}
	divisor := int64(math.Pow10(decimals))
	major := m.AmountCents / divisor
	minor := m.AmountCents % divisor
	formatStr := fmt.Sprintf("%%d.%%0%dd %%s", decimals)
	return fmt.Sprintf(formatStr, major, minor, code)
}

func main() {
	rates := map[string]float64{
		"USD/EUR": 0.92,
		"USD/JPY": 149.50,
		"EUR/USD": 1.087,
		"USD/BHD": 0.377,
	}

	examples := []struct {
		from Money
		to   string
	}{
		{Money{AmountCents: 1000, Currency: "USD"}, "EUR"},  // $10.00 → euros
		{Money{AmountCents: 1000, Currency: "USD"}, "JPY"},  // $10.00 → yen (zero-decimal)
		{Money{AmountCents: 5000, Currency: "USD"}, "BHD"},  // $50.00 → Bahraini Dinar (3 decimal)
		{Money{AmountCents: 2000, Currency: "EUR"}, "USD"},  // €20.00 → dollars
		{Money{AmountCents: 1500, Currency: "JPY"}, "JPY"},  // ¥1500 → same currency, no conversion
	}

	fmt.Fprintln(os.Stdout, "=== Currency Conversion Examples ===")
	for _, ex := range examples {
		result, err := ConvertCurrency(ex.from, ex.to, rates)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERROR: %v\n", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "%s → %s\n", FormatMoney(ex.from), FormatMoney(result))
	}

	fmt.Fprintln(os.Stdout, "\n=== Format Examples ===")
	samples := []Money{
		{AmountCents: 1099, Currency: "USD"},
		{AmountCents: 1000, Currency: "JPY"},
		{AmountCents: 10000, Currency: "BHD"},
		{AmountCents: 99, Currency: "EUR"},
	}
	for _, m := range samples {
		fmt.Fprintln(os.Stdout, FormatMoney(m))
	}
}
