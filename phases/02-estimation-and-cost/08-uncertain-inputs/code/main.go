package main

type Range struct {
	Low  float64
	Base float64
	High float64
}

func SpreadRatio(r Range) float64 {
	if r.Low == 0 {
		return 0
	}
	return r.High / r.Low
}

func IsWide(r Range, threshold float64) bool {
	return SpreadRatio(r) >= threshold
}

func main() {}
