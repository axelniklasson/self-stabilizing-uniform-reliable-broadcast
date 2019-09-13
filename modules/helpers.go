package modules

func isSubset(s, s2 map[int]bool) bool {
	for key := range s {
		if _, exists := s2[key]; !exists {
			return false
		}
	}
	return true
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func contains(set []int, val int) bool {
	for _, x := range set {
		if x == val {
			return true
		}
	}
	return false
}
