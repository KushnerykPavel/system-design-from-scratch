package main

import (
	"errors"
	"testing"
)

// --- IsTerminal ---

func TestIsTerminal_ReadIsTerminal(t *testing.T) {
	if !IsTerminal(READ) {
		t.Error("READ should be a terminal state")
	}
}

func TestIsTerminal_FailedIsTerminal(t *testing.T) {
	if !IsTerminal(FAILED) {
		t.Error("FAILED should be a terminal state")
	}
}

func TestIsTerminal_NonTerminalStates(t *testing.T) {
	for _, s := range []MessageState{QUEUED, SENT, DELIVERED} {
		if IsTerminal(s) {
			t.Errorf("%s should not be a terminal state", s)
		}
	}
}

// --- Valid transitions ---

func TestTransition_QueuedToSentOnServerAck(t *testing.T) {
	msg := &Message{ID: "t1", State: QUEUED}
	if err := Transition(msg, "server_ack"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.State != SENT {
		t.Errorf("expected SENT, got %s", msg.State)
	}
}

func TestTransition_SentToDeliveredOnDeviceAck(t *testing.T) {
	msg := &Message{ID: "t2", State: SENT}
	if err := Transition(msg, "device_ack"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.State != DELIVERED {
		t.Errorf("expected DELIVERED, got %s", msg.State)
	}
}

func TestTransition_DeliveredToReadOnReadAck(t *testing.T) {
	msg := &Message{ID: "t3", State: DELIVERED}
	if err := Transition(msg, "read_ack"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.State != READ {
		t.Errorf("expected READ, got %s", msg.State)
	}
}

func TestTransition_QueuedToFailedOnSendError(t *testing.T) {
	msg := &Message{ID: "t4", State: QUEUED}
	if err := Transition(msg, "send_error"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.State != FAILED {
		t.Errorf("expected FAILED, got %s", msg.State)
	}
}

func TestTransition_SentToFailedOnTTLExpiry(t *testing.T) {
	msg := &Message{ID: "t5", State: SENT}
	if err := Transition(msg, "ttl_expiry"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.State != FAILED {
		t.Errorf("expected FAILED, got %s", msg.State)
	}
}

// --- Full happy path ---

func TestTransition_FullDeliveryPath(t *testing.T) {
	msg := &Message{ID: "t6", State: QUEUED}
	events := []struct {
		event string
		want  MessageState
	}{
		{"server_ack", SENT},
		{"device_ack", DELIVERED},
		{"read_ack", READ},
	}
	for _, step := range events {
		if err := Transition(msg, step.event); err != nil {
			t.Fatalf("event %q failed: %v", step.event, err)
		}
		if msg.State != step.want {
			t.Errorf("after event %q: expected %s, got %s", step.event, step.want, msg.State)
		}
	}
	if !IsTerminal(msg.State) {
		t.Error("READ should be terminal after full delivery path")
	}
}

// --- Invalid transitions ---

func TestTransition_InvalidEventInState(t *testing.T) {
	msg := &Message{ID: "t7", State: QUEUED}
	err := Transition(msg, "read_ack") // read_ack is not valid from QUEUED
	if err == nil {
		t.Fatal("expected error for invalid event, got nil")
	}
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
	// State must not change on error
	if msg.State != QUEUED {
		t.Errorf("state should remain QUEUED after invalid event, got %s", msg.State)
	}
}

func TestTransition_SkippedState(t *testing.T) {
	// Attempt to go from SENT directly to READ, skipping DELIVERED
	msg := &Message{ID: "t8", State: SENT}
	err := Transition(msg, "read_ack")
	if err == nil {
		t.Fatal("expected error for read_ack from SENT, got nil")
	}
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
}

func TestTransition_TerminalStateRejectsAllEvents(t *testing.T) {
	for _, terminalState := range []MessageState{READ, FAILED} {
		for _, event := range []string{"server_ack", "device_ack", "read_ack", "send_error", "ttl_expiry"} {
			msg := &Message{ID: "t9", State: terminalState}
			err := Transition(msg, event)
			if err == nil {
				t.Errorf("expected error for event %q in terminal state %s, got nil", event, terminalState)
			}
			if !errors.Is(err, ErrInvalidTransition) {
				t.Errorf("expected ErrInvalidTransition, got %v", err)
			}
		}
	}
}

// --- Attempt counter ---

func TestTransition_AttemptsIncremented(t *testing.T) {
	msg := &Message{ID: "t10", State: QUEUED}
	_ = Transition(msg, "server_ack")
	_ = Transition(msg, "device_ack")
	if msg.Attempts != 2 {
		t.Errorf("expected 2 attempts after 2 successful transitions, got %d", msg.Attempts)
	}
}

func TestTransition_AttemptsNotIncrementedOnError(t *testing.T) {
	msg := &Message{ID: "t11", State: QUEUED}
	_ = Transition(msg, "invalid_event")
	if msg.Attempts != 0 {
		t.Errorf("attempts should not increment on failed transition, got %d", msg.Attempts)
	}
}
