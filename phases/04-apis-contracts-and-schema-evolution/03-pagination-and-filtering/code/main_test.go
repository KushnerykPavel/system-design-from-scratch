package main

import "testing"

func TestUsesCursor(t *testing.T) {
	if !UsesCursor(50, 1, true) {
		t.Fatal("mutable ordering should prefer cursor pagination")
	}
	if !UsesCursor(50, 1000, false) {
		t.Fatal("deep pages should prefer cursor pagination")
	}
	if UsesCursor(50, 2, false) != true {
		t.Fatal("default list endpoints should bias toward cursors")
	}
}
