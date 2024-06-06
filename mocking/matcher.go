package mocking

import "reflect"

// ArgumentMatcher can check whether an argument matches an expectation.
type ArgumentMatcher interface {
	// IsMatch returns true when the value of the argument matches the expectation.
	IsMatch(other any) bool

	// MatcherType returns the type of the matcher.
	MatcherType() MatcherType
}

// MatcherType defines the type of the matcher.
type MatcherType uint

const (
	// MatcherTypeFn is a matcher that uses a custom match function.
	MatcherTypeFn MatcherType = iota

	// MatcherTypeNone is used to specify that no arguments are passed to a variadic function.
	MatcherTypeNone

	// MatcherTypeAnyArgs is used to specify that any number of arguments can be passed to a variadic function.
	MatcherTypeAnyArgs

	// MatcherTypeVariadic is a matcher that matches the variable argument list passed to a variadic function.
	MatcherTypeVariadic
)

// Matcher is used to match an argument in a method invocation.
type Matcher[T any] struct {
	MatchFn     func(input T) bool
	matcherType MatcherType
}

// IsMatch returns true if other is a match to the expectation.
func (i Matcher[T]) IsMatch(other any) bool {
	return i.MatchFn(other.(T))
}

// MatcherType returns the matcher's type.
func (i Matcher[T]) MatcherType() MatcherType {
	return i.matcherType
}

type variadicMatcher struct {
	matchers []ArgumentMatcher
}

// Variadic creates a matcher for matching variable parameter lists. It accepts a slice of matchers
// that it delegates to for matching individual parameters.
func Variadic(matchers []ArgumentMatcher) ArgumentMatcher {
	return &variadicMatcher{matchers: matchers}
}

// MatcherType returns the type of the matcher.
func (v *variadicMatcher) MatcherType() MatcherType {
	return MatcherTypeVariadic
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

	if len(v.matchers) == 1 {
		switch v.matchers[0].MatcherType() {
		case MatcherTypeNone:
			return len(args) == 0
		case MatcherTypeAnyArgs:
			return true
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

// None is used to indicate that no arguments should be passed to a variadic function.
func None[T any]() Matcher[T] {
	return Matcher[T]{
		matcherType: MatcherTypeNone,
	}
}

// AnyArgs is used to indicate that any amount of arguments (including no arguments) should
// be passed to a variadic function.
func AnyArgs[T any]() Matcher[T] {
	return Matcher[T]{
		matcherType: MatcherTypeAnyArgs,
	}
}

// MethodMatcher is used to match a method call to an expectation.
type MethodMatcher struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
	Times            *uint
}
