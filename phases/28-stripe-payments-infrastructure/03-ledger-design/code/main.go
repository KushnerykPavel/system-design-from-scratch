package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Direction represents the direction of a ledger entry.
type Direction string

const (
	Debit  Direction = "DEBIT"
	Credit Direction = "CREDIT"
)

// Entry is a single immutable ledger entry.
type Entry struct {
	ID            string
	AccountID     string
	Currency      string
	AmountCents   int64
	Direction     Direction
	TransactionID string
}

// Transaction groups a set of double-entry ledger entries.
type Transaction struct {
	ID      string
	Entries []Entry
}

var (
	// ErrTooFewEntries is returned when a transaction has fewer than two entries.
	ErrTooFewEntries = errors.New("transaction must have at least 2 entries")
	// ErrZeroAmount is returned when an entry has a zero or negative amount.
	ErrZeroAmount = errors.New("entry amount must be greater than zero")
	// ErrUnbalanced is returned when debits and credits do not balance per currency.
	ErrUnbalanced = errors.New("transaction debits and credits are not balanced")
)

// ValidateTransaction checks that:
//   - the transaction has at least 2 entries
//   - no entry has a zero or negative amount
//   - for each currency, sum(debits) == sum(credits)
func ValidateTransaction(tx Transaction) error {
	if len(tx.Entries) < 2 {
		return ErrTooFewEntries
	}

	type currencyBalance struct {
		debits  int64
		credits int64
	}
	balances := make(map[string]*currencyBalance)

	for _, e := range tx.Entries {
		if e.AmountCents <= 0 {
			return fmt.Errorf("%w: entry %s has amount %d", ErrZeroAmount, e.ID, e.AmountCents)
		}
		if _, ok := balances[e.Currency]; !ok {
			balances[e.Currency] = &currencyBalance{}
		}
		switch e.Direction {
		case Debit:
			balances[e.Currency].debits += e.AmountCents
		case Credit:
			balances[e.Currency].credits += e.AmountCents
		}
	}

	for currency, b := range balances {
		if b.debits != b.credits {
			return fmt.Errorf("%w: currency %s debits=%d credits=%d",
				ErrUnbalanced, currency, b.debits, b.credits)
		}
	}

	return nil
}

// ComputeBalance returns the net balance (credits minus debits) for a given
// account and currency across all provided entries.
func ComputeBalance(accountID, currency string, entries []Entry) int64 {
	var balance int64
	for _, e := range entries {
		if e.AccountID != accountID || e.Currency != currency {
			continue
		}
		switch e.Direction {
		case Credit:
			balance += e.AmountCents
		case Debit:
			balance -= e.AmountCents
		}
	}
	return balance
}

func main() {
	// Valid transaction: merchant receives $50, Stripe takes $1.47 fee, net $48.53 to merchant.
	validTx := Transaction{
		ID: "txn_valid_001",
		Entries: []Entry{
			{ID: "e1", AccountID: "acct_customer", Currency: "USD", AmountCents: 5000, Direction: Debit, TransactionID: "txn_valid_001"},
			{ID: "e2", AccountID: "acct_merchant", Currency: "USD", AmountCents: 4853, Direction: Credit, TransactionID: "txn_valid_001"},
			{ID: "e3", AccountID: "acct_stripe_fees", Currency: "USD", AmountCents: 147, Direction: Credit, TransactionID: "txn_valid_001"},
		},
	}

	// Invalid transaction: credits exceed debits.
	invalidTx := Transaction{
		ID: "txn_invalid_001",
		Entries: []Entry{
			{ID: "e4", AccountID: "acct_customer", Currency: "USD", AmountCents: 1000, Direction: Debit, TransactionID: "txn_invalid_001"},
			{ID: "e5", AccountID: "acct_merchant", Currency: "USD", AmountCents: 1500, Direction: Credit, TransactionID: "txn_invalid_001"},
		},
	}

	// Validate transactions.
	validErr := ValidateTransaction(validTx)
	invalidErr := ValidateTransaction(invalidTx)

	// Compute balance for merchant account across both transactions.
	allEntries := append(validTx.Entries, invalidTx.Entries...)
	merchantBalance := ComputeBalance("acct_merchant", "USD", allEntries)

	result := map[string]interface{}{
		"valid_transaction_error":   fmt.Sprintf("%v", validErr),
		"invalid_transaction_error": fmt.Sprintf("%v", invalidErr),
		"merchant_balance_cents":    merchantBalance,
		"merchant_balance_usd":      fmt.Sprintf("$%.2f", float64(merchantBalance)/100),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
