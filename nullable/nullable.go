// Package nullable contains helpers for helping wrap and unwrap types with pointers.
package nullable

// OfValue returns a pointer to the specified value.
func OfValue[T any](value T) *T {
	return &value
}
