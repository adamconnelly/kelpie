// Code generated by Kelpie. DO NOT EDIT.
package interfacefilter

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

func (m *Instance) Include(name string) (r0 bool) {
	expectation := m.mock.Call("Include", name)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(name string) bool)
			return observe(name)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(bool)
		}
	}

	return
}

func (m *Mock) Instance() *Instance {
	return &m.instance
}

type IncludeMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *IncludeMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func Include[P0 string | mocking.Matcher[string]](name P0) *IncludeMethodMatcher {
	result := IncludeMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "Include",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(name).(mocking.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(name).(string))
	}

	return &result
}

type IncludeTimes struct {
	matcher *IncludeMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *IncludeMethodMatcher) Times(times uint) *IncludeTimes {
	m.matcher.Times = &times

	return &IncludeTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *IncludeMethodMatcher) Once() *IncludeTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *IncludeMethodMatcher) Never() *IncludeTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *IncludeTimes) Return(r0 bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *IncludeTimes) Panic(arg any) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *IncludeTimes) When(observe func(name string) bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *IncludeTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *IncludeMethodMatcher) Return(r0 bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *IncludeMethodMatcher) Panic(arg any) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *IncludeMethodMatcher) When(observe func(name string) bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type IncludeAction struct {
	expectation mocking.Expectation
}

func (a *IncludeAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
