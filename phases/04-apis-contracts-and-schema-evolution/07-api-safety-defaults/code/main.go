package main

func SafePageSize(requested, max int) int {
	if requested <= 0 {
		return min(max, 50)
	}
	if requested > max {
		return max
	}
	return requested
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {}
