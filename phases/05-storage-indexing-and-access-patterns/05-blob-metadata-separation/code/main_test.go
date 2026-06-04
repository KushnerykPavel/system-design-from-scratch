package main

import "testing"

func TestObjectStates(t *testing.T) {
	if !IsOrphaned(ObjectRecord{BlobPresent: true, MetadataPresent: false, State: StatePending}) {
		t.Fatal("expected object without metadata to be orphaned")
	}
	if CanServe(ObjectRecord{BlobPresent: true, MetadataPresent: true, State: StateDeleting}) {
		t.Fatal("expected deleting object not to be serveable")
	}
	if !CanServe(ObjectRecord{BlobPresent: true, MetadataPresent: true, State: StateReady}) {
		t.Fatal("expected ready object to be serveable")
	}
}
