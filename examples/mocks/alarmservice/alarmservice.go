// Code generated by Kelpie. DO NOT EDIT.
package alarmservice

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

func (m *Instance) CreateAlarm(name string) (r0 error) {
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

func (m *Mock) Instance() *Instance {
	return &m.instance
}

type CreateAlarmMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *CreateAlarmMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func CreateAlarm[P0 string | mocking.Matcher[string]](name P0) *CreateAlarmMethodMatcher {
	result := CreateAlarmMethodMatcher{
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

func (a *CreateAlarmMethodMatcher) Return(r0 error) *CreateAlarmAction {
	return &CreateAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{r0},
		},
	}
}

func (a *CreateAlarmMethodMatcher) Panic(arg any) *CreateAlarmAction {
	return &CreateAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *CreateAlarmMethodMatcher) When(observe func(name string) error) *CreateAlarmAction {
	return &CreateAlarmAction{
		expectation: mocking.Expectation{
			MethodMatcher: &a.matcher,
			ObserveFn:     observe,
		},
	}
}

type CreateAlarmAction struct {
	expectation mocking.Expectation
}

func (a *CreateAlarmAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}

func (a *CreateAlarmAction) Times(times int) *CreateAlarmTimes {
	a.expectation.MethodMatcher.Times = &times

	return &CreateAlarmTimes{
		expectation: a.expectation,
	}
}

func (a *CreateAlarmAction) Once() *CreateAlarmTimes {
	times := 1
	a.expectation.MethodMatcher.Times = &times

	return &CreateAlarmTimes{
		expectation: a.expectation,
	}
}

type CreateAlarmTimes struct {
	expectation mocking.Expectation
}

func (t *CreateAlarmTimes) CreateExpectation() *mocking.Expectation {
	return &t.expectation
}
