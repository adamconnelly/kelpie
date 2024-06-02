package mocking

// ArgumentMatcher can check whether an argument matches an expectation.
type ArgumentMatcher interface {
	// IsMatch returns true when the value of the argument matches the expectation.
	IsMatch(other any) bool
}

// Matcher is used to match an argument in a method invocation.
type Matcher[T any] struct {
	MatchFn func(input T) bool
}

// IsMatch returns true if other is a match to the expectation.
func (i Matcher[T]) IsMatch(other any) bool {
	return i.MatchFn(other.(T))
}

// MethodMatcher is used to match a method call to an expectation.
type MethodMatcher struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
	Times            *uint
}
