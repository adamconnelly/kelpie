// Code generated by Kelpie. DO NOT EDIT.
package maths

import "github.com/adamconnelly/kelpie"

type Mock struct {
	kelpie.Mock
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

func (m *Instance) Add(a int, b int) (r0 int) {
	expectation := m.mock.Call("Add", a, b)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(a int, b int) int)
			return observe(a, b)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(int)
		}
	}

	return
}

func (m *Instance) ParseInt(input string) (r0 int, r1 error) {
	expectation := m.mock.Call("ParseInt", input)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(input string) (int, error))
			return observe(input)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(int)
		}

		if expectation.Returns[1] != nil {
			r1 = expectation.Returns[1].(error)
		}
	}

	return
}

func (m *Mock) Instance() *Instance {
	return &m.instance
}

type AddMethodMatcher struct {
	matcher kelpie.MethodMatcher
}

func (m *AddMethodMatcher) CreateMethodMatcher() *kelpie.MethodMatcher {
	return &m.matcher
}

func Add[P0 int | kelpie.Matcher[int], P1 int | kelpie.Matcher[int]](a P0, b P1) *AddMethodMatcher {
	result := AddMethodMatcher{
		matcher: kelpie.MethodMatcher{
			MethodName:       "Add",
			ArgumentMatchers: make([]kelpie.ArgumentMatcher, 2),
		},
	}

	if matcher, ok := any(a).(kelpie.Matcher[int]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(a).(int))
	}

	if matcher, ok := any(b).(kelpie.Matcher[int]); ok {
		result.matcher.ArgumentMatchers[1] = matcher
	} else {
		result.matcher.ArgumentMatchers[1] = kelpie.ExactMatch(any(b).(int))
	}

	return &result
}

type AddExpectation struct {
	expectation kelpie.Expectation
}

func (e *AddExpectation) CreateExpectation() *kelpie.Expectation {
	return &e.expectation
}

func (a *AddMethodMatcher) Return(r0 int) *AddExpectation {
	return &AddExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{r0},
		},
	}
}

func (a *AddMethodMatcher) Panic(arg any) *AddExpectation {
	return &AddExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *AddMethodMatcher) When(observe func(a int, b int) int) *AddExpectation {
	return &AddExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			ObserveFn:     observe,
		},
	}
}

type ParseIntMethodMatcher struct {
	matcher kelpie.MethodMatcher
}

func (m *ParseIntMethodMatcher) CreateMethodMatcher() *kelpie.MethodMatcher {
	return &m.matcher
}

func ParseInt[P0 string | kelpie.Matcher[string]](input P0) *ParseIntMethodMatcher {
	result := ParseIntMethodMatcher{
		matcher: kelpie.MethodMatcher{
			MethodName:       "ParseInt",
			ArgumentMatchers: make([]kelpie.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(input).(kelpie.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(input).(string))
	}

	return &result
}

type ParseIntExpectation struct {
	expectation kelpie.Expectation
}

func (e *ParseIntExpectation) CreateExpectation() *kelpie.Expectation {
	return &e.expectation
}

func (a *ParseIntMethodMatcher) Return(r0 int, r1 error) *ParseIntExpectation {
	return &ParseIntExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{r0, r1},
		},
	}
}

func (a *ParseIntMethodMatcher) Panic(arg any) *ParseIntExpectation {
	return &ParseIntExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *ParseIntMethodMatcher) When(observe func(input string) (int, error)) *ParseIntExpectation {
	return &ParseIntExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			ObserveFn:     observe,
		},
	}
}
