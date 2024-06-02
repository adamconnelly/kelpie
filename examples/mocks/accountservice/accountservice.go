// Code generated by Kelpie. DO NOT EDIT.
package accountservice

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

func (m *Instance) SendActivationEmail(emailAddress string) (r0 bool) {
	expectation := m.mock.Call("SendActivationEmail", emailAddress)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(emailAddress string) bool)
			return observe(emailAddress)
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

func (m *Instance) DisableAccount(id uint) {
	expectation := m.mock.Call("DisableAccount", id)
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func(id uint))
			observe(id)
			return
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}
	}

	return
}

func (m *Instance) DisabledAccountIDs() (r0 []uint) {
	expectation := m.mock.Call("DisabledAccountIDs")
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func() []uint)
			return observe()
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].([]uint)
		}
	}

	return
}

func (m *Instance) DisableReasons() (r0 map[uint]string) {
	expectation := m.mock.Call("DisableReasons")
	if expectation != nil {
		if expectation.ObserveFn != nil {
			observe := expectation.ObserveFn.(func() map[uint]string)
			return observe()
		}

		if expectation.PanicArg != nil {
			panic(expectation.PanicArg)
		}

		if expectation.Returns[0] != nil {
			r0 = expectation.Returns[0].(map[uint]string)
		}
	}

	return
}

func (m *Mock) Instance() *Instance {
	return &m.instance
}

type SendActivationEmailMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *SendActivationEmailMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func SendActivationEmail[P0 string | mocking.Matcher[string]](emailAddress P0) *SendActivationEmailMethodMatcher {
	result := SendActivationEmailMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "SendActivationEmail",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(emailAddress).(mocking.Matcher[string]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(emailAddress).(string))
	}

	return &result
}

type SendActivationEmailTimes struct {
	matcher *SendActivationEmailMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *SendActivationEmailMethodMatcher) Times(times uint) *SendActivationEmailTimes {
	m.matcher.Times = &times

	return &SendActivationEmailTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *SendActivationEmailMethodMatcher) Once() *SendActivationEmailTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *SendActivationEmailMethodMatcher) Never() *SendActivationEmailTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *SendActivationEmailTimes) Return(r0 bool) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *SendActivationEmailTimes) Panic(arg any) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *SendActivationEmailTimes) When(observe func(emailAddress string) bool) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *SendActivationEmailTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *SendActivationEmailMethodMatcher) Return(r0 bool) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *SendActivationEmailMethodMatcher) Panic(arg any) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *SendActivationEmailMethodMatcher) When(observe func(emailAddress string) bool) *SendActivationEmailAction {
	return &SendActivationEmailAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type SendActivationEmailAction struct {
	expectation mocking.Expectation
}

func (a *SendActivationEmailAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}

type DisableAccountMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *DisableAccountMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func DisableAccount[P0 uint | mocking.Matcher[uint]](id P0) *DisableAccountMethodMatcher {
	result := DisableAccountMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "DisableAccount",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 1),
		},
	}

	if matcher, ok := any(id).(mocking.Matcher[uint]); ok {
		result.matcher.ArgumentMatchers[0] = matcher
	} else {
		result.matcher.ArgumentMatchers[0] = kelpie.ExactMatch(any(id).(uint))
	}

	return &result
}

type DisableAccountTimes struct {
	matcher *DisableAccountMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *DisableAccountMethodMatcher) Times(times uint) *DisableAccountTimes {
	m.matcher.Times = &times

	return &DisableAccountTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *DisableAccountMethodMatcher) Once() *DisableAccountTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *DisableAccountMethodMatcher) Never() *DisableAccountTimes {
	return m.Times(0)
}

// Panic panics using the specified argument when the method is called.
func (t *DisableAccountTimes) Panic(arg any) *DisableAccountAction {
	return &DisableAccountAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *DisableAccountTimes) When(observe func(id uint)) *DisableAccountAction {
	return &DisableAccountAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *DisableAccountTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Panic panics using the specified argument when the method is called.
func (m *DisableAccountMethodMatcher) Panic(arg any) *DisableAccountAction {
	return &DisableAccountAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *DisableAccountMethodMatcher) When(observe func(id uint)) *DisableAccountAction {
	return &DisableAccountAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type DisableAccountAction struct {
	expectation mocking.Expectation
}

func (a *DisableAccountAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}

type DisabledAccountIDsMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *DisabledAccountIDsMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func DisabledAccountIDs() *DisabledAccountIDsMethodMatcher {
	result := DisabledAccountIDsMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "DisabledAccountIDs",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 0),
		},
	}

	return &result
}

type DisabledAccountIDsTimes struct {
	matcher *DisabledAccountIDsMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *DisabledAccountIDsMethodMatcher) Times(times uint) *DisabledAccountIDsTimes {
	m.matcher.Times = &times

	return &DisabledAccountIDsTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *DisabledAccountIDsMethodMatcher) Once() *DisabledAccountIDsTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *DisabledAccountIDsMethodMatcher) Never() *DisabledAccountIDsTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *DisabledAccountIDsTimes) Return(r0 []uint) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *DisabledAccountIDsTimes) Panic(arg any) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *DisabledAccountIDsTimes) When(observe func() []uint) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *DisabledAccountIDsTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *DisabledAccountIDsMethodMatcher) Return(r0 []uint) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *DisabledAccountIDsMethodMatcher) Panic(arg any) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *DisabledAccountIDsMethodMatcher) When(observe func() []uint) *DisabledAccountIDsAction {
	return &DisabledAccountIDsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type DisabledAccountIDsAction struct {
	expectation mocking.Expectation
}

func (a *DisabledAccountIDsAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}

type DisableReasonsMethodMatcher struct {
	matcher mocking.MethodMatcher
}

func (m *DisableReasonsMethodMatcher) CreateMethodMatcher() *mocking.MethodMatcher {
	return &m.matcher
}

func DisableReasons() *DisableReasonsMethodMatcher {
	result := DisableReasonsMethodMatcher{
		matcher: mocking.MethodMatcher{
			MethodName:       "DisableReasons",
			ArgumentMatchers: make([]mocking.ArgumentMatcher, 0),
		},
	}

	return &result
}

type DisableReasonsTimes struct {
	matcher *DisableReasonsMethodMatcher
}

// Times allows you to restrict the number of times a particular expectation can be matched.
func (m *DisableReasonsMethodMatcher) Times(times uint) *DisableReasonsTimes {
	m.matcher.Times = &times

	return &DisableReasonsTimes{
		matcher: m,
	}
}

// Once specifies that the expectation will only match once.
func (m *DisableReasonsMethodMatcher) Once() *DisableReasonsTimes {
	return m.Times(1)
}

// Never specifies that the method has not been called. This is mainly useful for verification
// rather than mocking.
func (m *DisableReasonsMethodMatcher) Never() *DisableReasonsTimes {
	return m.Times(0)
}

// Return returns the specified results when the method is called.
func (t *DisableReasonsTimes) Return(r0 map[uint]string) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (t *DisableReasonsTimes) Panic(arg any) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (t *DisableReasonsTimes) When(observe func() map[uint]string) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &t.matcher.matcher,
			ObserveFn:     observe,
		},
	}
}

func (t *DisableReasonsTimes) CreateMethodMatcher() *mocking.MethodMatcher {
	return &t.matcher.matcher
}

// Return returns the specified results when the method is called.
func (m *DisableReasonsMethodMatcher) Return(r0 map[uint]string) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			Returns:       []any{r0},
		},
	}
}

// Panic panics using the specified argument when the method is called.
func (m *DisableReasonsMethodMatcher) Panic(arg any) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			PanicArg:      arg,
		},
	}
}

// When calls the specified observe callback when the method is called.
func (m *DisableReasonsMethodMatcher) When(observe func() map[uint]string) *DisableReasonsAction {
	return &DisableReasonsAction{
		expectation: mocking.Expectation{
			MethodMatcher: &m.matcher,
			ObserveFn:     observe,
		},
	}
}

type DisableReasonsAction struct {
	expectation mocking.Expectation
}

func (a *DisableReasonsAction) CreateExpectation() *mocking.Expectation {
	return &a.expectation
}
