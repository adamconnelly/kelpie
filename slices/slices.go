package slices

// FirstOrPanic returns the first item matching the closure, and panics if there are no matches.
func FirstOrPanic[T any](slice []T, matches func(item T) bool) (result T) {
	for _, item := range slice {
		if matches(item) {
			return item
		}
	}

	panic("Item not found in slice")
}

func Contains[T any](slice []T, matches func(item T) bool) bool {
	for _, item := range slice {
		if matches(item) {
			return true
		}
	}

	return false
}
