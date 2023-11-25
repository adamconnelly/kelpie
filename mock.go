package kelpie

import "fmt"

type Expectation struct {
	MethodName       string
	ArgumentMatchers []ArgumentMatcher
	Returns          []any
	PanicArg         any
	ObserveFn        any
}

type MethodCall struct {
	MethodName string
	Args       []any
}

type Mock struct {
	Expectations []*Expectation
	MethodCalls  []*MethodCall
}

func (m *Mock) Setup(expectation *Expectation) {
	m.Expectations = append([]*Expectation{expectation}, m.Expectations...)
}

func (m *Mock) Call(methodName string, args ...any) *Expectation {
	m.MethodCalls = append(m.MethodCalls, &MethodCall{MethodName: methodName, Args: args})

	for _, expectation := range m.Expectations {
		if expectation.MethodName == methodName {
			if len(args) != len(expectation.ArgumentMatchers) {
				panic(fmt.Sprintf("Argument mismatch in call to '%s'.\n    Expected: %d\n    Actual: %d\nThis is a bug in Kelpie - please report it!", methodName, len(expectation.ArgumentMatchers), len(args)))
			}

			argsMatch := true
			for i, matcher := range expectation.ArgumentMatchers {
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

func (m *Mock) Reset() {
	m.Expectations = nil
	m.MethodCalls = nil
}
