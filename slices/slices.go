package slices

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}
