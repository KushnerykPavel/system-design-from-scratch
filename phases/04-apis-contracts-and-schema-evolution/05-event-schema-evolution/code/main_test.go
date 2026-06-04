package main

import "testing"

func TestCompatible(t *testing.T) {
	if !Compatible(AddOptional) {
		t.Fatal("optional additions should be compatible by default")
	}
	if Compatible(AddRequired) {
		t.Fatal("required additions should not be compatible")
	}
	if Compatible(RenameField) {
		t.Fatal("renames should not be treated as compatible")
	}
}
