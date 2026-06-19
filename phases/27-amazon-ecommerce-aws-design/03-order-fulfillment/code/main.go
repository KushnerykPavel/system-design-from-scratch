package main

import (
	"errors"
	"fmt"
)

// OrderState represents the lifecycle state of an order.
type OrderState string

const (
	StatePlaced            OrderState = "PLACED"
	StatePaymentPending    OrderState = "PAYMENT_PENDING"
	StatePaymentConfirmed  OrderState = "PAYMENT_CONFIRMED"
	StateInventoryReserved OrderState = "INVENTORY_RESERVED"
	StatePicking           OrderState = "PICKING"
	StateShipped           OrderState = "SHIPPED"
	StateDelivered         OrderState = "DELIVERED"
	StateCancelled         OrderState = "CANCELLED"
)

// validTransitions defines which state transitions are legal.
var validTransitions = map[OrderState][]OrderState{
	StatePlaced: {
		StatePaymentPending,
		StateCancelled,
	},
	StatePaymentPending: {
		StatePaymentConfirmed,
		StateCancelled,
	},
	StatePaymentConfirmed: {
		StateInventoryReserved,
		StateCancelled,
	},
	StateInventoryReserved: {
		StatePicking,
		StateCancelled,
	},
	StatePicking: {
		StateShipped,
	},
	StateShipped: {
		StateDelivered,
	},
	// Terminal states — no outbound transitions.
	StateDelivered: {},
	StateCancelled: {},
}

// Order represents a customer order.
type Order struct {
	ID             string
	State          OrderState
	Items          []string
	IdempotencyKey string
}

// NewOrder creates an order in the PLACED state.
func NewOrder(id string, items []string, idempotencyKey string) *Order {
	return &Order{
		ID:             id,
		State:          StatePlaced,
		Items:          items,
		IdempotencyKey: idempotencyKey,
	}
}

// ErrInvalidTransition is returned when a state transition is not allowed.
var ErrInvalidTransition = errors.New("invalid state transition")

// Transition moves the order to newState if the transition is valid.
// Applying the same transition twice (idempotent re-try) returns nil.
func Transition(order *Order, newState OrderState) error {
	// Idempotent: already in target state is a no-op.
	if order.State == newState {
		return nil
	}

	allowed, ok := validTransitions[order.State]
	if !ok {
		return fmt.Errorf("%w: unknown source state %q", ErrInvalidTransition, order.State)
	}
	for _, s := range allowed {
		if s == newState {
			order.State = newState
			return nil
		}
	}
	return fmt.Errorf("%w: %q → %q", ErrInvalidTransition, order.State, newState)
}

func main() {
	fmt.Println("=== Happy path: PLACED → DELIVERED ===")
	order := NewOrder("ORD-001", []string{"WIDGET-A", "WIDGET-B"}, "idem-001")
	happyPath := []OrderState{
		StatePaymentPending,
		StatePaymentConfirmed,
		StateInventoryReserved,
		StatePicking,
		StateShipped,
		StateDelivered,
	}
	for _, s := range happyPath {
		if err := Transition(order, s); err != nil {
			fmt.Printf("  ERROR: %v\n", err)
			return
		}
		fmt.Printf("  → %s\n", order.State)
	}

	fmt.Println("\n=== Cancellation path: PLACED → CANCELLED ===")
	order2 := NewOrder("ORD-002", []string{"GADGET-X"}, "idem-002")
	if err := Transition(order2, StateCancelled); err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  → %s\n", order2.State)
	}

	fmt.Println("\n=== Cancellation after payment: PAYMENT_CONFIRMED → CANCELLED ===")
	order3 := NewOrder("ORD-003", []string{"CAMERA-1"}, "idem-003")
	steps := []OrderState{StatePaymentPending, StatePaymentConfirmed}
	for _, s := range steps {
		_ = Transition(order3, s)
	}
	fmt.Printf("  Current state: %s\n", order3.State)
	if err := Transition(order3, StateCancelled); err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  → %s (compensation: refund triggered)\n", order3.State)
	}

	fmt.Println("\n=== Invalid: SHIPPED → CANCELLED (should fail) ===")
	order4 := NewOrder("ORD-004", []string{"BOOK-1"}, "idem-004")
	for _, s := range []OrderState{StatePaymentPending, StatePaymentConfirmed, StateInventoryReserved, StatePicking, StateShipped} {
		_ = Transition(order4, s)
	}
	if err := Transition(order4, StateCancelled); err != nil {
		fmt.Printf("  Correctly rejected: %v\n", err)
	} else {
		fmt.Printf("  ERROR: should have been rejected\n")
	}

	fmt.Println("\n=== Idempotent re-try: applying same transition twice ===")
	order5 := NewOrder("ORD-005", []string{"PHONE-1"}, "idem-005")
	_ = Transition(order5, StatePaymentPending)
	if err := Transition(order5, StatePaymentPending); err != nil {
		fmt.Printf("  ERROR: idempotency failed: %v\n", err)
	} else {
		fmt.Printf("  Idempotent re-apply accepted, state still: %s\n", order5.State)
	}
}
