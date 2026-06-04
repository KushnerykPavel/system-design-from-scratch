package main

type FieldChange string

const (
	AddOptional FieldChange = "add_optional"
	AddRequired FieldChange = "add_required"
	RenameField FieldChange = "rename_field"
)

func Compatible(change FieldChange) bool {
	return change == AddOptional
}

func main() {}
