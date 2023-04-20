package utils

func SliceEquals[T comparable](t1 []T, t2 []T) bool {
	if len(t1) != len(t2) {
		return false
	}

	for i, t := range t1 {
		if t2[i] != t {
			return false
		}
	}

	return true
}
