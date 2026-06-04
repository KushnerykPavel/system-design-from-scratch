package main

type ChangeType string

const (
	AddOptionalField   ChangeType = "add_optional_field"
	RemoveField        ChangeType = "remove_field"
	ChangeFieldMeaning ChangeType = "change_field_meaning"
)

func IsBreaking(change ChangeType) bool {
	switch change {
	case AddOptionalField:
		return false
	case RemoveField, ChangeFieldMeaning:
		return true
	default:
		return true
	}
}

func main() {}
