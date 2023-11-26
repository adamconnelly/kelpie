// Code generated by Kelpie. DO NOT EDIT.
package interfacefilter

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

func (m *Instance) Include(name string) (r0 bool) {
	expectation := m.mock.Call("Include", name)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(name string) (bool))
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
	matcher kelpie.MethodMatcher
}

func (m *IncludeMethodMatcher) CreateMethodMatcher() *kelpie.MethodMatcher {
	return &m.matcher
}

func Include[P0 string | kelpie.Matcher[string]](name P0) *IncludeMethodMatcher {
	result := IncludeMethodMatcher{
		matcher: kelpie.MethodMatcher{
			MethodName: "Include",
			ArgumentMatchers: make([]kelpie.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(name).(kelpie.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(name).(string))
	}

	return &result
}

type IncludeExpectation struct {
	expectation kelpie.Expectation
}

func (e *IncludeExpectation) CreateExpectation() *kelpie.Expectation {
	return &e.expectation
}

func (a *IncludeMethodMatcher) Return(r0 bool) *IncludeExpectation {
	return &IncludeExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{r0},
		},
	}
}

func (a *IncludeMethodMatcher) Panic(arg any) *IncludeExpectation {
	return &IncludeExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *IncludeMethodMatcher) When(observe func(name string) (bool)) *IncludeExpectation {
	return &IncludeExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			ObserveFn:     observe,
		},
	}
}
