package main

func UsesCursor(pageSize, pageNumber int, mutableOrder bool) bool {
	if mutableOrder {
		return true
	}
	if pageNumber > 100 {
		return true
	}
	return pageSize > 0
}

func main() {}
