package main

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreakerOpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker("test-svc", 3, 30*time.Second)
	fail := func() error { return errors.New("fail") }
	for i := 0; i < 3; i++ {
		_ = cb.Call(fail)
	}
	if State(cb.state) != StateOpen {
		t.Fatalf("expected circuit to be open after 3 failures, got %s", State(cb.state).String())
	}
}

func TestCircuitBreakerBlocksCallWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker("test-svc", 1, 30*time.Second)
	_ = cb.Call(func() error { return errors.New("fail") })
	err := cb.Call(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerResetsOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker("test-svc", 5, 30*time.Second)
	fail := func() error { return errors.New("fail") }
	for i := 0; i < 2; i++ {
		_ = cb.Call(fail)
	}
	// A success should reset failure count and keep circuit closed.
	err := cb.Call(func() error { return nil })
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if State(cb.state) != StateClosed {
		t.Fatalf("expected closed after success, got %s", State(cb.state).String())
	}
}
