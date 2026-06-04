package main

type Record struct {
	RequestHash string
	Status      string
}

type Decision string

const (
	Execute  Decision = "execute"
	Replay   Decision = "replay"
	Conflict Decision = "conflict"
)

func Decide(existing *Record, requestHash string) Decision {
	if existing == nil {
		return Execute
	}
	if existing.RequestHash != requestHash {
		return Conflict
	}
	return Replay
}

func main() {}
