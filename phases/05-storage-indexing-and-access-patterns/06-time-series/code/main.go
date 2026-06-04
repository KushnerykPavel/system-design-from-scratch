package main

func RecommendBucketSeconds(windowSeconds int) int {
	switch {
	case windowSeconds <= 3600:
		return 60
	case windowSeconds <= 86400:
		return 300
	default:
		return 3600
	}
}

func main() {}
