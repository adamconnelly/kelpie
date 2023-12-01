package mocking

import (
	"fmt"

	"github.com/adamconnelly/kelpie/slices"
)

type MethodMatcher struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
	Times            *uint
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

		if expectation.MethodMatcher.Times != nil {
			calls := slices.All(m.MethodCalls, func(methodCall *MethodCall) bool {
				return methodMatchesExpectation(methodMatcher, methodCall.MethodName, methodCall.Args...)
			})
			if uint(len(calls)) <= *expectation.MethodMatcher.Times {
				return expectation
			}
		} else if methodMatchesExpectation(methodMatcher, methodName, args...) {
			return expectation
		}
	}

	return nil
}

func (m *Mock) Called(creator MethodMatcherCreator) bool {
	methodMatcher := creator.CreateMethodMatcher()

	if methodMatcher.Times != nil {
		matches := slices.All(m.MethodCalls, func(call *MethodCall) bool {
			return methodMatchesExpectation(methodMatcher, call.MethodName, call.Args...)
		})

		return uint(len(matches)) == *methodMatcher.Times
	}

	return slices.Contains(m.MethodCalls, func(call *MethodCall) bool {
		return methodMatchesExpectation(methodMatcher, call.MethodName, call.Args...)
	})
}

func (m *Mock) Reset() {
	m.Expectations = nil
	m.MethodCalls = nil
}

func methodMatchesExpectation(methodMatcher *MethodMatcher, methodName string, args ...any) bool {
	if methodMatcher.MethodName == methodName {
		if len(args) != len(methodMatcher.ArgumentMatchers) {
			panic(fmt.Sprintf("Argument mismatch in call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", methodMatcher.MethodName, len(methodMatcher.ArgumentMatchers), len(args)))
		}

		argsMatch := true
		for i, matcher := range methodMatcher.ArgumentMatchers {
			if !matcher.IsMatch(args[i]) {
				argsMatch = false
				break
			}
		}

		if argsMatch {
			return true
		}
	}

	return false
}
