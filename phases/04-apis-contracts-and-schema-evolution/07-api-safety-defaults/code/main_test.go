package main

import "testing"

func TestSafePageSize(t *testing.T) {
	if got := SafePageSize(0, 100); got != 50 {
		t.Fatalf("expected default 50, got %d", got)
	}
	if got := SafePageSize(500, 100); got != 100 {
		t.Fatalf("expected cap 100, got %d", got)
	}
	if got := SafePageSize(20, 100); got != 20 {
		t.Fatalf("expected requested size preserved, got %d", got)
	}
}
