package main

type InterfaceKind string

const (
	HTTP  InterfaceKind = "http"
	GRPC  InterfaceKind = "grpc"
	Event InterfaceKind = "event"
)

type Workload struct {
	ExternalClients     bool
	NeedsImmediateReply bool
	HighFanout          bool
	ReplayNeeded        bool
	LowLatencyInternal  bool
}

func ChooseInterface(w Workload) InterfaceKind {
	if w.HighFanout || w.ReplayNeeded {
		return Event
	}
	if w.ExternalClients {
		return HTTP
	}
	if w.NeedsImmediateReply && w.LowLatencyInternal {
		return GRPC
	}
	if w.NeedsImmediateReply {
		return HTTP
	}
	return Event
}

func main() {}
