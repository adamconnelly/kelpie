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

// Contains returns true if any of the elements match using the supplied function.
func Contains[T any](slice []T, matches func(item T) bool) bool {
	for _, item := range slice {
		if matches(item) {
			return true
		}
	}

	return false
}

// All returns all the elements that match using the supplied function.
func All[T any](slice []T, matches func(item T) bool) []T {
	var results []T
	for _, item := range slice {
		if matches(item) {
			results = append(results, item)
		}
	}

	return results
}

// Map takes a source slice, and converts the elements to a destination type using the supplied map function.
func Map[TSource any, TDest any](slice []TSource, mapFn func(item TSource) TDest) []TDest {
	result := make([]TDest, len(slice))

	for i := range slice {
		result[i] = mapFn(slice[i])
	}

	return result
}
