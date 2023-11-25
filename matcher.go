package kelpie

// ArgumentMatcher can check whether an argument matches an expectation.
type ArgumentMatcher interface {
	// IsMatch returns true when the value of the argument matches the expectation.
	IsMatch(other any) bool
}

// Matcher is used to match an argument in a method invocation.
type Matcher[T comparable] struct {
	isMatch    func(input T) bool
	exactMatch T
	isAny      bool
}

// ExactMatch matches parameters to the value of exactMatch.
func ExactMatch[T comparable](exactMatch T) Matcher[T] {
	return Matcher[T]{exactMatch: exactMatch}
}

// IsMatch returns true if other is a match to the expectation.
func (i Matcher[T]) IsMatch(other any) bool {
	return (i.isMatch != nil && i.isMatch(other.(T))) || i.exactMatch == other.(T) || i.isAny
}

// Any matches any value of T.
func Any[T comparable]() Matcher[T] {
	return Matcher[T]{isAny: true}
}

// Match uses isMatch to determine whether an argument matches.
func Match[T comparable](isMatch func(arg T) bool) Matcher[T] {
	return Matcher[T]{isMatch: isMatch}
}
