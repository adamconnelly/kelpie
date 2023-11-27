package mocking

import (
	"fmt"
)

type MethodMatcher struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
}

type Expectation struct {
	MethodMatcher *MethodMatcher
	Returns       []any
	PanicArg      any
	ObserveFn     any
}

type MethodMatcherCreator interface {
	CreateMethodMatcher() *MethodMatcher
}

type ExpectationCreator interface {
	CreateExpectation() *Expectation
}

type MethodCall struct {
	MethodName string
	Args       []any
}

type Mock struct {
	Expectations []*Expectation
	MethodCalls  []*MethodCall
}

func (m *Mock) Setup(creator ExpectationCreator) {
	m.Expectations = append([]*Expectation{creator.CreateExpectation()}, m.Expectations...)
}

func (m *Mock) Call(methodName string, args ...any) *Expectation {
	m.MethodCalls = append(m.MethodCalls, &MethodCall{MethodName: methodName, Args: args})

	for _, expectation := range m.Expectations {
		methodMatcher := expectation.MethodMatcher
		if methodMatcher.MethodName == methodName {
			if len(args) != len(methodMatcher.ArgumentMatchers) {
				panic(fmt.Sprintf("Argument mismatch in call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", methodName, len(methodMatcher.ArgumentMatchers), len(args)))
			}

			argsMatch := true
			for i, matcher := range methodMatcher.ArgumentMatchers {
				if !matcher.IsMatch(args[i]) {
					argsMatch = false
					break
				}
			}

			if argsMatch {
				return expectation
			}
		}
	}

	return nil
}

func (m *Mock) Called(creator MethodMatcherCreator) bool {
	methodMatcher := creator.CreateMethodMatcher()
	for _, methodCall := range m.MethodCalls {
		if methodMatcher.MethodName == methodCall.MethodName {
			if len(methodCall.Args) != len(methodMatcher.ArgumentMatchers) {
				panic(fmt.Sprintf("Argument mismatch when checking call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", methodMatcher.MethodName, len(methodMatcher.ArgumentMatchers), len(methodCall.Args)))
			}

			argsMatch := true
			for i, matcher := range methodMatcher.ArgumentMatchers {
				if !matcher.IsMatch(methodCall.Args[i]) {
					argsMatch = false
					break
				}
			}

			if argsMatch {
				return true
			}
		}
	}

	return false
}

func (m *Mock) Reset() {
	m.Expectations = nil
	m.MethodCalls = nil
}
