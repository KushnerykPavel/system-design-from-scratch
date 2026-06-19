package main

import (
	"errors"
	"fmt"
)

// MessageState represents a stage in the message delivery lifecycle.
type MessageState int

const (
	// QUEUED: message is waiting to be sent (e.g., device is offline or pending network).
	QUEUED MessageState = iota
	// SENT: server has persisted the message and acked receipt to the sender.
	SENT
	// DELIVERED: recipient device has received the message and sent a delivery ack.
	DELIVERED
	// READ: recipient has opened the conversation; message was read.
	READ
	// FAILED: delivery failed after retry exhaustion or TTL expiry.
	FAILED
)

func (s MessageState) String() string {
	switch s {
	case QUEUED:
		return "QUEUED"
	case SENT:
		return "SENT"
	case DELIVERED:
		return "DELIVERED"
	case READ:
		return "READ"
	case FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// Message represents a message in the delivery pipeline.
type Message struct {
	ID       string
	State    MessageState
	Attempts int
}

// validTransitions defines all legal state transitions.
// A transition is identified by the event string received in a given state.
var validTransitions = map[MessageState]map[string]MessageState{
	QUEUED: {
		"server_ack": SENT,
		"send_error": FAILED,
	},
	SENT: {
		"device_ack": DELIVERED,
		"ttl_expiry": FAILED,
	},
	DELIVERED: {
		"read_ack": READ,
	},
	// READ and FAILED are terminal — no transitions out of them.
}

// ErrInvalidTransition is returned when an event cannot be applied to the current state.
var ErrInvalidTransition = errors.New("invalid state transition")

// Transition applies an event to a message, advancing its state machine.
// It returns ErrInvalidTransition if the event is not valid for the current state.
// Terminal states (READ, FAILED) reject all events.
func Transition(msg *Message, event string) error {
	if IsTerminal(msg.State) {
		return fmt.Errorf("%w: message %s is in terminal state %s", ErrInvalidTransition, msg.ID, msg.State)
	}
	transitions, ok := validTransitions[msg.State]
	if !ok {
		return fmt.Errorf("%w: no transitions defined from state %s", ErrInvalidTransition, msg.State)
	}
	nextState, ok := transitions[event]
	if !ok {
		return fmt.Errorf("%w: event %q not valid in state %s", ErrInvalidTransition, event, msg.State)
	}
	msg.State = nextState
	msg.Attempts++
	return nil
}

// IsTerminal returns true if the state is a terminal state (READ or FAILED).
// Terminal states cannot transition further.
func IsTerminal(state MessageState) bool {
	return state == READ || state == FAILED
}

func main() {
	messages := []struct {
		id     string
		events []string
	}{
		{
			id:     "msg-001",
			events: []string{"server_ack", "device_ack", "read_ack"},
		},
		{
			id:     "msg-002",
			events: []string{"server_ack", "ttl_expiry"},
		},
		{
			id:     "msg-003",
			events: []string{"send_error"},
		},
		{
			id:     "msg-004",
			events: []string{"server_ack", "device_ack"},
		},
		{
			id:     "msg-005",
			events: []string{"server_ack", "read_ack"}, // invalid: skip DELIVERED
		},
	}

	for _, tc := range messages {
		msg := &Message{ID: tc.id, State: QUEUED}
		fmt.Printf("Message %s:\n", msg.ID)
		fmt.Printf("  initial state: %s\n", msg.State)

		for _, event := range tc.events {
			prev := msg.State
			err := Transition(msg, event)
			if err != nil {
				fmt.Printf("  event %-12s -> ERROR: %v\n", event, err)
			} else {
				fmt.Printf("  event %-12s -> %s -> %s\n", event, prev, msg.State)
			}
		}

		fmt.Printf("  final state: %s (terminal=%v, attempts=%d)\n\n", msg.State, IsTerminal(msg.State), msg.Attempts)
	}
}
