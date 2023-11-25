package kelpie

import "fmt"

// TODO: Create MethodSetup, and MethodConfiguration. Setup will just contain the matching details, config will contain return args, etc.

type Expectation interface {
	MethodName() string
	ArgumentMatchers() []ArgumentMatcher
	Returns() []any
	PanicArg() any
	ObserveFn() any
}

type E struct {
	Method  string
	Args    []ArgumentMatcher
	Ret     []any
	Panic   any
	Observe any
}

func (e *E) MethodName() string {
	return e.Method
}

func (e *E) ArgumentMatchers() []ArgumentMatcher {
	return e.Args
}

func (e *E) Returns() []any {
	return e.Ret
}

func (e *E) PanicArg() any {
	return e.Panic
}

func (e *E) ObserveFn() any {
	return e.Observe
}

type MethodCall struct {
	MethodName string
	Args       []any
}

type Mock struct {
	Expectations []Expectation
	MethodCalls  []*MethodCall
}

// TODO: create a separate interface to use for Setup vs Called
func (m *Mock) Setup(expectation Expectation) {
	m.Expectations = append([]Expectation{expectation}, m.Expectations...)
}

func (m *Mock) Call(methodName string, args ...any) Expectation {
	m.MethodCalls = append(m.MethodCalls, &MethodCall{MethodName: methodName, Args: args})

	for _, expectation := range m.Expectations {
		if expectation.MethodName() == methodName {
			if len(args) != len(expectation.ArgumentMatchers()) {
				panic(fmt.Sprintf("Argument mismatch in call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", methodName, len(expectation.ArgumentMatchers()), len(args)))
			}

			argsMatch := true
			for i, matcher := range expectation.ArgumentMatchers() {
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

func (m *Mock) Called(expectation Expectation) bool {
	for _, methodCall := range m.MethodCalls {
		if expectation.MethodName() == methodCall.MethodName {
			if len(methodCall.Args) != len(expectation.ArgumentMatchers()) {
				panic(fmt.Sprintf("Argument mismatch when checking call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", expectation.MethodName(), len(expectation.ArgumentMatchers()), len(methodCall.Args)))
			}

			argsMatch := true
			for i, matcher := range expectation.ArgumentMatchers() {
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
