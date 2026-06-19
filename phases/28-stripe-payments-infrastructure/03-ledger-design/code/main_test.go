package main

import (
	"errors"
	"testing"
)

func TestValidateTransactionBalanced(t *testing.T) {
	tx := Transaction{
		ID: "txn_balanced",
		Entries: []Entry{
			{ID: "e1", AccountID: "acct_a", Currency: "USD", AmountCents: 1000, Direction: Debit, TransactionID: "txn_balanced"},
			{ID: "e2", AccountID: "acct_b", Currency: "USD", AmountCents: 1000, Direction: Credit, TransactionID: "txn_balanced"},
		},
	}
	if err := ValidateTransaction(tx); err != nil {
		t.Fatalf("expected valid transaction, got error: %v", err)
	}
}

func TestValidateTransactionUnbalanced(t *testing.T) {
	tx := Transaction{
		ID: "txn_unbalanced",
		Entries: []Entry{
			{ID: "e1", AccountID: "acct_a", Currency: "USD", AmountCents: 1000, Direction: Debit, TransactionID: "txn_unbalanced"},
			{ID: "e2", AccountID: "acct_b", Currency: "USD", AmountCents: 500, Direction: Credit, TransactionID: "txn_unbalanced"},
		},
	}
	err := ValidateTransaction(tx)
	if !errors.Is(err, ErrUnbalanced) {
		t.Fatalf("expected ErrUnbalanced, got %v", err)
	}
}

func TestValidateTransactionTooFewEntries(t *testing.T) {
	tx := Transaction{
		ID: "txn_one_entry",
		Entries: []Entry{
			{ID: "e1", AccountID: "acct_a", Currency: "USD", AmountCents: 1000, Direction: Debit, TransactionID: "txn_one_entry"},
		},
	}
	err := ValidateTransaction(tx)
	if !errors.Is(err, ErrTooFewEntries) {
		t.Fatalf("expected ErrTooFewEntries, got %v", err)
	}
}

func TestValidateTransactionZeroAmount(t *testing.T) {
	tx := Transaction{
		ID: "txn_zero",
		Entries: []Entry{
			{ID: "e1", AccountID: "acct_a", Currency: "USD", AmountCents: 0, Direction: Debit, TransactionID: "txn_zero"},
			{ID: "e2", AccountID: "acct_b", Currency: "USD", AmountCents: 0, Direction: Credit, TransactionID: "txn_zero"},
		},
	}
	err := ValidateTransaction(tx)
	if !errors.Is(err, ErrZeroAmount) {
		t.Fatalf("expected ErrZeroAmount, got %v", err)
	}
}

func TestComputeBalanceSingleTransaction(t *testing.T) {
	entries := []Entry{
		{ID: "e1", AccountID: "acct_merchant", Currency: "USD", AmountCents: 5000, Direction: Credit, TransactionID: "txn1"},
		{ID: "e2", AccountID: "acct_customer", Currency: "USD", AmountCents: 5000, Direction: Debit, TransactionID: "txn1"},
	}
	balance := ComputeBalance("acct_merchant", "USD", entries)
	if balance != 5000 {
		t.Fatalf("expected 5000, got %d", balance)
	}
}

func TestComputeBalanceMultipleTransactions(t *testing.T) {
	entries := []Entry{
		// Transaction 1: +$50 credit to merchant
		{ID: "e1", AccountID: "acct_merchant", Currency: "USD", AmountCents: 5000, Direction: Credit, TransactionID: "txn1"},
		// Transaction 2: $10 refund debits merchant
		{ID: "e2", AccountID: "acct_merchant", Currency: "USD", AmountCents: 1000, Direction: Debit, TransactionID: "txn2"},
	}
	balance := ComputeBalance("acct_merchant", "USD", entries)
	if balance != 4000 {
		t.Fatalf("expected 4000 (50-10 dollars), got %d", balance)
	}
}

func TestComputeBalanceCurrencyIsolation(t *testing.T) {
	entries := []Entry{
		{ID: "e1", AccountID: "acct_a", Currency: "USD", AmountCents: 10000, Direction: Credit, TransactionID: "txn1"},
		{ID: "e2", AccountID: "acct_a", Currency: "EUR", AmountCents: 8000, Direction: Credit, TransactionID: "txn2"},
	}
	usdBalance := ComputeBalance("acct_a", "USD", entries)
	eurBalance := ComputeBalance("acct_a", "EUR", entries)
	if usdBalance != 10000 {
		t.Fatalf("expected USD 10000, got %d", usdBalance)
	}
	if eurBalance != 8000 {
		t.Fatalf("expected EUR 8000, got %d", eurBalance)
	}
}
