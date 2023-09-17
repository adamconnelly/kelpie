// Code generated by Kelpie. DO NOT EDIT.
package interfacefilter

import "github.com/adamconnelly/kelpie"

type Mock struct {
	expectations []Expectation
}

func (m *Mock) Include(name string) (r0 bool) {
    for _, expectation := range m.expectations {
		if expectation.method == "Include" {
			info := expectation.invocationDetails.(IncludeInvocationDetails)
            if info.name.IsMatch(name) {
				if info.observe != nil {
					return info.observe(name)
				}

				if info.panicArg != nil {
					panic(info.panicArg)
				}

                return info.result0
			}
		}
	}

    return
}

type Expectation struct {
	method            string
	invocationDetails interface{}
}

func (m *Mock) Reset() {
	m.expectations = nil
}

func (m *Mock) Setup(expectation Expectation) {
	m.expectations = append([]Expectation{expectation}, m.expectations...)
}


type IncludeInvocationDetails struct {
    name kelpie.Matcher[string]
    result0 bool
	panicArg any
	observe func(name string) (bool)
}

func Include[P0 string | kelpie.Matcher[string]](name P0) IncludeInvocationDetails {
	result := IncludeInvocationDetails{}

    if matcher, ok := any(name).(kelpie.Matcher[string]); ok {
        result.name = matcher
    } else {
        result.name = kelpie.ExactMatch(any(name).(string))
    }

    return result
}

func (a IncludeInvocationDetails) Return(r0 bool) Expectation {
    a.result0 = r0

	return Expectation{
		method:            "Include",
		invocationDetails: a,
	}
}

func (a IncludeInvocationDetails) Panic(arg any) Expectation {
	a.panicArg = arg

	return Expectation{
		method:            "Include",
		invocationDetails: a,
	}
}

func (a IncludeInvocationDetails) When(callback func(name string) (bool)) Expectation {
	a.observe = callback

	return Expectation{
		method:            "Include",
		invocationDetails: a,
	}
}
