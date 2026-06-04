package main

import "testing"

func TestRequiresContractGate(t *testing.T) {
	if !RequiresContractGate(Change{ShapeChanged: true}) {
		t.Fatal("shape changes should require contract checks")
	}
	if !RequiresContractGate(Change{SemanticsChanged: true}) {
		t.Fatal("semantic changes should require contract checks")
	}
	if RequiresContractGate(Change{}) {
		t.Fatal("no contract-relevant change should not require a gate")
	}
}
