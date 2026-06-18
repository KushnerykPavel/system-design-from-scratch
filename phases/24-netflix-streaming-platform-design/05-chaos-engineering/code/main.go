package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// State represents a circuit breaker state.
type State int

const (
	StateClosed   State = iota // normal: calls pass through
	StateOpen                  // tripped: calls fail fast
	StateHalfOpen              // testing: one probe call allowed
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreaker implements a simple error-rate circuit breaker.
type CircuitBreaker struct {
	name         string
	threshold    int // consecutive failures before opening
	resetTimeout time.Duration
	failures     int32
	state        int32 // atomic State
	openedAt     int64 // unix nano
}

// NewCircuitBreaker creates a circuit breaker with the given failure threshold.
func NewCircuitBreaker(name string, threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

// Call executes fn, tracking failures and state transitions.
func (cb *CircuitBreaker) Call(fn func() error) error {
	s := State(atomic.LoadInt32(&cb.state))
	if s == StateOpen {
		// Check if reset timeout has elapsed.
		opened := time.Unix(0, atomic.LoadInt64(&cb.openedAt))
		if time.Since(opened) > cb.resetTimeout {
			atomic.StoreInt32(&cb.state, int32(StateHalfOpen))
		} else {
			return ErrCircuitOpen
		}
	}
	err := fn()
	if err != nil {
		failures := atomic.AddInt32(&cb.failures, 1)
		if int(failures) >= cb.threshold {
			atomic.StoreInt32(&cb.state, int32(StateOpen))
			atomic.StoreInt64(&cb.openedAt, time.Now().UnixNano())
		}
		return err
	}
	// Success: reset.
	atomic.StoreInt32(&cb.failures, 0)
	atomic.StoreInt32(&cb.state, int32(StateClosed))
	return nil
}

// Status returns a snapshot of the current breaker state.
func (cb *CircuitBreaker) Status() map[string]any {
	return map[string]any{
		"name":     cb.name,
		"state":    State(atomic.LoadInt32(&cb.state)).String(),
		"failures": atomic.LoadInt32(&cb.failures),
	}
}

func main() {
	cb := NewCircuitBreaker("recommendation-service", 3, 10*time.Second)

	// Simulate 3 consecutive failures -> circuit opens.
	alwaysFail := func() error { return errors.New("upstream error") }
	for i := 0; i < 3; i++ {
		_ = cb.Call(alwaysFail)
	}

	// 4th call should be rejected by the open circuit.
	err := cb.Call(alwaysFail)

	result := map[string]any{
		"status":              cb.Status(),
		"fourth_call_blocked": errors.Is(err, ErrCircuitOpen),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if encErr := enc.Encode(result); encErr != nil {
		fmt.Fprintln(os.Stderr, encErr)
		os.Exit(1)
	}
}
