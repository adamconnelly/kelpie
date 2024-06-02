// Package kelpie contains helpers for matching arguments when configuring a mock.
package kelpie

import (
	"reflect"

	"github.com/adamconnelly/kelpie/mocking"
)

// ExactMatch matches parameters to the value of exactMatch.
func ExactMatch[T any](exactMatch T) mocking.Matcher[T] {
	return Match(func(arg T) bool {
		return reflect.DeepEqual(arg, exactMatch)
	})
}

// Any matches any value of T.
func Any[T any]() mocking.Matcher[T] {
	return Match(func(arg T) bool { return true })
}

// Match uses isMatch to determine whether an argument matches.
func Match[T any](isMatch func(arg T) bool) mocking.Matcher[T] {
	return mocking.Matcher[T]{MatchFn: isMatch}
}
