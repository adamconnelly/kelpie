package kelpie

import "github.com/adamconnelly/kelpie/mocking"

// ExactMatch matches parameters to the value of exactMatch.
func ExactMatch[T comparable](exactMatch T) mocking.Matcher[T] {
	return Match(func(arg T) bool { return arg == exactMatch })
}

// Any matches any value of T.
func Any[T comparable]() mocking.Matcher[T] {
	return Match(func(arg T) bool { return true })
}

// Match uses isMatch to determine whether an argument matches.
func Match[T comparable](isMatch func(arg T) bool) mocking.Matcher[T] {
	return mocking.Matcher[T]{MatchFn: isMatch}
}
