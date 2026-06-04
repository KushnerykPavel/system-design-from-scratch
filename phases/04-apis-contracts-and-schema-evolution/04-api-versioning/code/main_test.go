package main

import "testing"

func TestIsBreaking(t *testing.T) {
	if IsBreaking(AddOptionalField) {
		t.Fatal("adding an optional field should be compatible by default")
	}
	if !IsBreaking(RemoveField) {
		t.Fatal("removing a field should be breaking")
	}
	if !IsBreaking(ChangeFieldMeaning) {
		t.Fatal("semantic changes should be treated as breaking")
	}
}
