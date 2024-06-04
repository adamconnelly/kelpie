package mocking

import "reflect"

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

type variadicMatcher struct {
	matchers []ArgumentMatcher
}

// Variadic creates a matcher for matching variable parameter lists. It accepts a slice of matchers
// that it delegates to for matching individual parameters.
func Variadic(matchers []ArgumentMatcher) ArgumentMatcher {
	return &variadicMatcher{matchers: matchers}
}

// IsMatch returns true if other is a match to the expectation.
func (v *variadicMatcher) IsMatch(other any) bool {
	args, ok := other.([]any)
	if !ok {
		// If the arguments can't be cast to []any, they might still be a slice of a different
		// type. In this case we can use reflection to copy the arguments into a []any so
		// that we can compare them.
		t := reflect.TypeOf(other)
		if t.Kind() != reflect.Slice {
			return false
		}

		value := reflect.ValueOf(other)
		args = make([]any, value.Len())

		for i := 0; i < value.Len(); i++ {
			args[i] = value.Index(i).Interface()
		}
	}

	if len(args) != len(v.matchers) {
		return false
	}

	for index, matcher := range v.matchers {
		if !matcher.IsMatch(args[index]) {
			return false
		}
	}

	return true
}

// MethodMatcher is used to match a method call to an expectation.
type MethodMatcher struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
	Times            *uint
}
