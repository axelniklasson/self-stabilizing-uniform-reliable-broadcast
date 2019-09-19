package ssurb

func isSubset(s, s2 map[int]bool) bool {
	for key := range s {
		if _, exists := s2[key]; !exists {
			return false
		}
	}
	return true
}

func listToMap(s []int) map[int]bool {
	m := map[int]bool{}
	for _, x := range s {
		m[x] = true
	}
	return m
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func contains(s []int, val int) bool {
	for _, x := range s {
		if x == val {
			return true
		}
	}
	return false
}
