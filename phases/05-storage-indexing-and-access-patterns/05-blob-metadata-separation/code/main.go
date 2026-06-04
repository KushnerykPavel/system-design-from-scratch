package main

type ObjectState string

const (
	StatePending  ObjectState = "pending"
	StateReady    ObjectState = "ready"
	StateDeleting ObjectState = "deleting"
)

type ObjectRecord struct {
	BlobPresent     bool
	MetadataPresent bool
	State           ObjectState
}

func IsOrphaned(record ObjectRecord) bool {
	return record.BlobPresent && !record.MetadataPresent
}

func CanServe(record ObjectRecord) bool {
	return record.BlobPresent && record.MetadataPresent && record.State == StateReady
}

func main() {}
