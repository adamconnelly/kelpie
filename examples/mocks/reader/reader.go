// Code generated by Kelpie. DO NOT EDIT.
package reader

import (
	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocking"
)

type Mock struct {
	mocking.Mock
	instance Instance
}

func NewMock() *Mock {
	mock := Mock{
		instance: Instance{},
	}
	mock.instance.mock = &mock

	return &mock
}

type Instance struct {
	mock *Mock
}

func (m *Instance) Read(p []byte) (n int, err error) {
	expectation := m.mock.Call("Read", p)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(p []byte) (int, error))
			return observe(p)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			n = expectation.Returns[0].(int)
		}

		if expectation.Returns[1] != nil {
			err = expectation.Returns[1].(error)
		}
	}

	return
}

func (m *Mock) Instance() *Instance {
	return &m.instance
}

type ReadMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *ReadMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func Read[P0 []byte | mocking.Matcher[[]byte]](p P0) *ReadMethodMatcher {
	result := ReadMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "Read",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(p).(mocking.Matcher[[]byte]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(p).([]byte))
	}

	return &result
}

type ReadTimes struct {
	matcher *ReadMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *ReadMethodMatcher) Times(times uint) *ReadTimes {
	m.matcher.Times = &times

	return &ReadTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *ReadMethodMatcher) Once() *ReadTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *ReadMethodMatcher) Never() *ReadTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *ReadTimes) Return(n int, err error) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{n, err},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *ReadTimes) Panic(arg any) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *ReadTimes) When(observe func(p []byte) (int, error)) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *ReadTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *ReadMethodMatcher) Return(n int, err error) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{n, err},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *ReadMethodMatcher) Panic(arg any) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *ReadMethodMatcher) When(observe func(p []byte) (int, error)) *ReadAction {
	return &ReadAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type ReadAction struct {
	expectation mocking.Expectation
}

func (a *ReadAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
