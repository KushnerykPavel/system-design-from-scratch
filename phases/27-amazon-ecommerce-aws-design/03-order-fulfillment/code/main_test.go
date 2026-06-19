package main

import (
	"errors"
	"testing"
)

func TestHappyPathTransitions(t *testing.T) {
	order := NewOrder("test-001", []string{"ITEM-A"}, "idem-001")
	path := []OrderState{
		StatePaymentPending,
		StatePaymentConfirmed,
		StateInventoryReserved,
		StatePicking,
		StateShipped,
		StateDelivered,
	}
	for _, s := range path {
		if err := Transition(order, s); err != nil {
			t.Fatalf("unexpected error transitioning to %q: %v", s, err)
		}
		if order.State != s {
			t.Fatalf("expected state %q, got %q", s, order.State)
		}
	}
}

func TestCancellationFromPlaced(t *testing.T) {
	order := NewOrder("test-002", []string{"ITEM-B"}, "idem-002")
	if err := Transition(order, StateCancelled); err != nil {
		t.Fatalf("expected PLACED → CANCELLED to succeed, got: %v", err)
	}
	if order.State != StateCancelled {
		t.Fatalf("expected CANCELLED, got %q", order.State)
	}
}

func TestCancellationFromPaymentPending(t *testing.T) {
	order := NewOrder("test-003", []string{"ITEM-C"}, "idem-003")
	_ = Transition(order, StatePaymentPending)
	if err := Transition(order, StateCancelled); err != nil {
		t.Fatalf("expected PAYMENT_PENDING → CANCELLED to succeed, got: %v", err)
	}
}

func TestCancellationFromPaymentConfirmed(t *testing.T) {
	order := NewOrder("test-004", []string{"ITEM-D"}, "idem-004")
	_ = Transition(order, StatePaymentPending)
	_ = Transition(order, StatePaymentConfirmed)
	if err := Transition(order, StateCancelled); err != nil {
		t.Fatalf("expected PAYMENT_CONFIRMED → CANCELLED to succeed, got: %v", err)
	}
}

func TestCancellationFromInventoryReserved(t *testing.T) {
	order := NewOrder("test-005", []string{"ITEM-E"}, "idem-005")
	_ = Transition(order, StatePaymentPending)
	_ = Transition(order, StatePaymentConfirmed)
	_ = Transition(order, StateInventoryReserved)
	if err := Transition(order, StateCancelled); err != nil {
		t.Fatalf("expected INVENTORY_RESERVED → CANCELLED to succeed, got: %v", err)
	}
}

func TestInvalidTransitionShippedToCancelled(t *testing.T) {
	order := NewOrder("test-006", []string{"ITEM-F"}, "idem-006")
	for _, s := range []OrderState{StatePaymentPending, StatePaymentConfirmed, StateInventoryReserved, StatePicking, StateShipped} {
		_ = Transition(order, s)
	}
	err := Transition(order, StateCancelled)
	if err == nil {
		t.Fatal("expected error for SHIPPED → CANCELLED, got nil")
	}
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

func TestInvalidTransitionDeliveredToAnyForwardState(t *testing.T) {
	order := NewOrder("test-007", []string{"ITEM-G"}, "idem-007")
	for _, s := range []OrderState{StatePaymentPending, StatePaymentConfirmed, StateInventoryReserved, StatePicking, StateShipped, StateDelivered} {
		_ = Transition(order, s)
	}
	// DELIVERED is terminal — no forward transitions allowed.
	for _, next := range []OrderState{StateCancelled, StateShipped, StatePicking} {
		err := Transition(order, next)
		if err == nil {
			t.Errorf("expected error for DELIVERED → %q, got nil", next)
		}
	}
}

func TestIdempotentTransition(t *testing.T) {
	order := NewOrder("test-008", []string{"ITEM-H"}, "idem-008")
	_ = Transition(order, StatePaymentPending)
	// Applying the same transition again must not error.
	if err := Transition(order, StatePaymentPending); err != nil {
		t.Fatalf("idempotent re-apply should return nil, got: %v", err)
	}
	if order.State != StatePaymentPending {
		t.Fatalf("state should remain PAYMENT_PENDING, got %q", order.State)
	}
}

func TestInvalidTransitionSkipsState(t *testing.T) {
	order := NewOrder("test-009", []string{"ITEM-I"}, "idem-009")
	// Try to jump directly from PLACED to INVENTORY_RESERVED (skipping payment steps).
	err := Transition(order, StateInventoryReserved)
	if err == nil {
		t.Fatal("expected error for skipping states, got nil")
	}
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

func TestCancelledIsTerminal(t *testing.T) {
	order := NewOrder("test-010", []string{"ITEM-J"}, "idem-010")
	_ = Transition(order, StateCancelled)
	// Cannot transition out of CANCELLED to any other state.
	err := Transition(order, StatePlaced)
	if err == nil {
		t.Fatal("expected error transitioning out of CANCELLED, got nil")
	}
}

func TestNewOrderStartsInPlaced(t *testing.T) {
	order := NewOrder("test-011", []string{"ITEM-K"}, "idem-011")
	if order.State != StatePlaced {
		t.Fatalf("expected initial state PLACED, got %q", order.State)
	}
}
