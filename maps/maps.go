// Package maps contains utilities for working with maps.
package maps

// Keys returns the unique set of keys in the map.
func Keys[TKey comparable, TValue any](m map[TKey]TValue) []TKey {
	keys := make([]TKey, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index++
	}

	return keys
}

// Values returns the values in the map.
func Values[TKey comparable, TValue any](m map[TKey]TValue) []TValue {
	values := make([]TValue, len(m))
	index := 0
	for _, v := range m {
		values[index] = v
		index++
	}

	return values
}
