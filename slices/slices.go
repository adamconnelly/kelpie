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

// Map takes a source slice, and converts the elements to a destination type using the supplied map function.
func Map[TSource any, TDest any](slice []TSource, mapFn func(item TSource) TDest) []TDest {
	result := make([]TDest, len(slice))

	for i := range slice {
		result[i] = mapFn(slice[i])
	}

	return result
}
