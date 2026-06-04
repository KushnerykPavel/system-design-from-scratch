package main

type Change struct {
	ShapeChanged     bool
	SemanticsChanged bool
}

func RequiresContractGate(change Change) bool {
	return change.ShapeChanged || change.SemanticsChanged
}

func main() {}
