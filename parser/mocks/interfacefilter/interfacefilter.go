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

func (a *IncludeMethodMatcher) Return(r0 bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{r0},
		},
	}
}

func (a *IncludeMethodMatcher) Panic(arg any) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *IncludeMethodMatcher) When(observe func(name string) bool) *IncludeAction {
	return &IncludeAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
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

func (a *IncludeAction) Times(times int) *IncludeTimes {
	a.expectation.MethodMatcher.Times = &times

	return &IncludeTimes{
		expectation: a.expectation,
	}
}

func (a *IncludeAction) Once() *IncludeTimes {
	times := 1
	a.expectation.MethodMatcher.Times = &times

	return &IncludeTimes{
		expectation: a.expectation,
	}
}

type IncludeTimes struct {
	expectation mocking.Expectation
}

func (t *IncludeTimes) CreateExpectation() *mocking.Expectation {
	return &t.expectation
}
