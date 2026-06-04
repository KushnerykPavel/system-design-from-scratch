package main

type BottleneckModel struct {
	CPUReqPerSecond     float64
	DiskReqPerSecond    float64
	NetworkReqPerSecond float64
}

func Bottleneck(model BottleneckModel) (string, float64) {
	name := "cpu"
	limit := model.CPUReqPerSecond

	if model.DiskReqPerSecond < limit {
		name = "disk"
		limit = model.DiskReqPerSecond
	}
	if model.NetworkReqPerSecond < limit {
		name = "network"
		limit = model.NetworkReqPerSecond
	}

	return name, limit
}

func main() {}
