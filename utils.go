package main

func divisionBoundaries(length int, parts int) [][]int {
	if parts <= 0 || length <= 0 {
		return [][]int{}
	}
	
	offsets := make([][]int, parts)
	base := length / parts
	start := 0

	for i := range parts {
		end := min(start + base, length)

		offsets[i] = []int{start, end - 1}
		start = end
	}
	return offsets
}
