// Code generated by Kelpie. DO NOT EDIT.
package alarmservice

import (
	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocking"
)

type Mock struct {
	mocking.Mock
	instance instance
}

func NewMock() *Mock {
	mock := Mock{
		instance: instance{},
	}
	mock.instance.mock = &mock

	return &mock
}

type instance struct {
	mock *Mock
}

func (m *instance) CreateAlarm(name string) (r0 error) {
	expectation := m.mock.Call("CreateAlarm", name)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(name string) error)
			return observe(name)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(error)
		}
	}

	return
}

func (m *Mock) Instance() *instance {
	return &m.instance
}

type createAlarmMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *createAlarmMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func CreateAlarm[P0 string | mocking.Matcher[string]](name P0) *createAlarmMethodMatcher {
	result := createAlarmMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "CreateAlarm",
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

type createAlarmTimes struct {
	matcher *createAlarmMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *createAlarmMethodMatcher) Times(times uint) *createAlarmTimes {
	m.matcher.Times = &times

	return &createAlarmTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *createAlarmMethodMatcher) Once() *createAlarmTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *createAlarmMethodMatcher) Never() *createAlarmTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *createAlarmTimes) Return(r0 error) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *createAlarmTimes) Panic(arg any) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *createAlarmTimes) When(observe func(name string) error) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *createAlarmTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *createAlarmMethodMatcher) Return(r0 error) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *createAlarmMethodMatcher) Panic(arg any) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *createAlarmMethodMatcher) When(observe func(name string) error) *createAlarmAction {
	return &createAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type createAlarmAction struct {
	expectation mocking.Expectation
}

func (a *createAlarmAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
