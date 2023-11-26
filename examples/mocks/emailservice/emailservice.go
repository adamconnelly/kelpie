// Code generated by Kelpie. DO NOT EDIT.
package emailservice

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

func (m *Instance) Send(sender string, recipient string, body string) (cost float64, err error) {
	expectation := m.mock.Call("Send", sender, recipient, body)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(sender string, recipient string, body string) (float64, error))
			return observe(sender, recipient, body)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		
		if expectation.Returns[0] != nil {
			cost = expectation.Returns[0].(float64)
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


type SendMethodMatcher struct {
	matcher kelpie.MethodMatcher
}

func (m *SendMethodMatcher) CreateMethodMatcher() *kelpie.MethodMatcher {
	return &m.matcher
}

func Send[P0 string | kelpie.Matcher[string], P1 string | kelpie.Matcher[string], P2 string | kelpie.Matcher[string]](sender P0, recipient P1, body P2) *SendMethodMatcher {
	result := SendMethodMatcher{
		matcher: kelpie.MethodMatcher{
			MethodName: "Send",
			ArgumentMatchers: make([]kelpie.ArgumentMatcher, 3),
		},
	}

	if matcher, ok := any(sender).(kelpie.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(sender).(string))
	}

	if matcher, ok := any(recipient).(kelpie.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[1] = matcher
	} else {
		result.matcher.ArgumentMatchers[1] = kelpie.ExactMatch(any(recipient).(string))
	}

	if matcher, ok := any(body).(kelpie.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[2] = matcher
	} else {
		result.matcher.ArgumentMatchers[2] = kelpie.ExactMatch(any(body).(string))
	}

	return &result
}

type SendExpectation struct {
	expectation kelpie.Expectation
}

func (e *SendExpectation) CreateExpectation() *kelpie.Expectation {
	return &e.expectation
}

func (a *SendMethodMatcher) Return(cost float64, err error) *SendExpectation {
	return &SendExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			Returns:       []any{cost, err},
		},
	}
}

func (a *SendMethodMatcher) Panic(arg any) *SendExpectation {
	return &SendExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			PanicArg:      arg,
		},
	}
}

func (a *SendMethodMatcher) When(observe func(sender string, recipient string, body string) (float64, error)) *SendExpectation {
	return &SendExpectation{
		expectation: kelpie.Expectation{
			MethodMatcher: &a.matcher,
			ObserveFn:     observe,
		},
	}
}