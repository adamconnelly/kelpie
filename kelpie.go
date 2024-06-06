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

// None is used when mocking methods that contain a variable parameter list to indicate that
// no parameters should be provided.
//
// For example: printMock.Setup(print.Printf("Testing 123", kelpie.None[any]())).
func None[T any]() mocking.Matcher[T] {
	return mocking.None[T]()
}

// AnyArgs is used when mocking methods that contain a variable parameter list when you don't
// care what arguments are passed as the variable parameter. For example:
//
// printMock.Setup(print.Printf("Testing 123", kelpie.AnyArgs[any]()).Panic("Oh no!"))
func AnyArgs[T any]() mocking.Matcher[T] {
	return mocking.AnyArgs[T]()
}
