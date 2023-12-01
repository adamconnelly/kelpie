// Package mocking contains the mocking logic for Kelpie. This includes types for creating
// expectations as well as the Mock type which is used by the generated mocks to record method
// calls and verify expectations.
package mocking

import (
	"fmt"

	"github.com/adamconnelly/kelpie/slices"
)

// Expectation represents an expected method call.
type Expectation struct {
	// MethodMatcher contains the information needed to match a specific method call to an expectation.
	MethodMatcher *MethodMatcher

	// Returns contains any arguments that should be returned from the method call.
	Returns []any

	// PanicArg contains the argument passed to `panic()`.
	PanicArg any

	// ObserveFn contains a function that should be used as the implementation of the mocked method.
	ObserveFn any
}

// MethodMatcherCreator is used to create a MethodMatcher. This is used to allow us to build
// the fluent API that ensures that only the details of an expected method call can be passed
// to `Called`, rather than allowing a full setup expression containing an action.
type MethodMatcherCreator interface {
	CreateMethodMatcher() *MethodMatcher
}

// ExpectationCreator creates an expected method call including the details needed to match
// the method, as well as the result that should occur (i.e. return something, panic or call
// a custom function).
type ExpectationCreator interface {
	CreateExpectation() *Expectation
}

// MethodCall is used to record a method that has been called.
type MethodCall struct {
	// MethodName is the name of the method that will be called.
	MethodName string

	// Args contains the list of arguments passed to the method.
	Args []any
}

// Mock contains the main mocking logic.
type Mock struct {
	// Expectations contains the list of expected method calls that have been setup.
	Expectations []*Expectation

	// MethodCalls contains the method calls that have been recorded.
	MethodCalls []*MethodCall
}

// Setup is used to configure a method call for a mock. The setup contains the information
// needed to match a specific method call, as well as the action that should occur when the
// method is called.
//
// Examples:
//
//	mock.Setup(calculator.Add(1, kelpie.Any[int]).Return(5))
//	mock.Setup(client.Request("abc").Times(3).Return(errors.New("request failed")))
func (m *Mock) Setup(creator ExpectationCreator) {
	m.Expectations = append([]*Expectation{creator.CreateExpectation()}, m.Expectations...)
}

// Call records a method call, and returns an expectation if any can be found. If no expectations
// match the specified method call, nil will be returned.
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

// Called verifies whether a method matching the specified signature has been called.
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

// Reset clears the expectations and recorded method calls on the mock.
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
