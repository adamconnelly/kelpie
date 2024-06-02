// Code generated by Kelpie. DO NOT EDIT.
package requester

import (
	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocking"
	"io"
	. "net/http"
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

func (m *instance) MakeRequest(r *Request) (r0 io.Reader, r1 error) {
	expectation := m.mock.Call("MakeRequest", r)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(r *Request) (io.Reader, error))
			return observe(r)
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(io.Reader)
		}

		if expectation.Returns[1] != nil {
			r1 = expectation.Returns[1].(error)
		}
	}

	return
}

func (m *Mock) Instance() *instance {
	return &m.instance
}

type makeRequestMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *makeRequestMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func MakeRequest[P0 *Request | mocking.Matcher[*Request]](r P0) *makeRequestMethodMatcher {
	result := makeRequestMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "MakeRequest",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(r).(mocking.Matcher[*Request]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(r).(*Request))
	}

	return &result
}

type makeRequestTimes struct {
	matcher *makeRequestMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *makeRequestMethodMatcher) Times(times uint) *makeRequestTimes {
	m.matcher.Times = &times

	return &makeRequestTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *makeRequestMethodMatcher) Once() *makeRequestTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *makeRequestMethodMatcher) Never() *makeRequestTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *makeRequestTimes) Return(r0 io.Reader, r1 error) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0, r1},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *makeRequestTimes) Panic(arg any) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *makeRequestTimes) When(observe func(r *Request) (io.Reader, error)) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *makeRequestTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *makeRequestMethodMatcher) Return(r0 io.Reader, r1 error) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0, r1},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *makeRequestMethodMatcher) Panic(arg any) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *makeRequestMethodMatcher) When(observe func(r *Request) (io.Reader, error)) *makeRequestAction {
	return &makeRequestAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type makeRequestAction struct {
	expectation mocking.Expectation
}

func (a *makeRequestAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
