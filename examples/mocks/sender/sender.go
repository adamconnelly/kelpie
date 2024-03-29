// Code generated by Kelpie. DO NOT EDIT.
package sender

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

func (m *Instance) SendMessage(title *string, message string) (r0 error) {
	expectation := m.mock.Call("SendMessage", title, message)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(title *string, message string) error)
			return observe(title, message)
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

type SendMessageMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *SendMessageMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func SendMessage[P0 *string | mocking.Matcher[*string], P1 string | mocking.Matcher[string]](title P0, message P1) *SendMessageMethodMatcher {
	result := SendMessageMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "SendMessage",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 2),
		},
	}

	if matcher, ok := any(title).(mocking.Matcher[*string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(title).(*string))
	}

	if matcher, ok := any(message).(mocking.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[1] = matcher
	} else {
		result.matcher.ArgumentMatchers[1] = kelpie.ExactMatch(any(message).(string))
	}

	return &result
}

type SendMessageTimes struct {
	matcher *SendMessageMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *SendMessageMethodMatcher) Times(times uint) *SendMessageTimes {
	m.matcher.Times = &times

	return &SendMessageTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *SendMessageMethodMatcher) Once() *SendMessageTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *SendMessageMethodMatcher) Never() *SendMessageTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *SendMessageTimes) Return(r0 error) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *SendMessageTimes) Panic(arg any) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *SendMessageTimes) When(observe func(title *string, message string) error) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *SendMessageTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *SendMessageMethodMatcher) Return(r0 error) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *SendMessageMethodMatcher) Panic(arg any) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *SendMessageMethodMatcher) When(observe func(title *string, message string) error) *SendMessageAction {
	return &SendMessageAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type SendMessageAction struct {
	expectation mocking.Expectation
}

func (a *SendMessageAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
