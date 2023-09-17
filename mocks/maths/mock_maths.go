package maths

import "github.com/adamconnelly/kelpie"

type Mock struct {
	expectations []Expectation
}

func (m *Mock) Reset() {
	m.expectations = nil
}

func (m *Mock) Add(a, b int) int {
	var result int

	for _, expectation := range m.expectations {
		if expectation.method == "Add" {
			info := expectation.info.(AddMethodCallBegin)
			if info.a.IsMatch(a) && info.b.IsMatch(b) {
				if info.callback != nil {
					return info.callback(a, b)
				}

				if info.panicArg != nil {
					panic(info.panicArg)
				}

				return info.result
			}
		}
	}

	return result
}

func (m *Mock) ParseInt(input string) (int, error) {
	var return0 int
	var return1 error

	for _, expectation := range m.expectations {
		if expectation.method == "ParseInt" {
			info := expectation.info.(ParseIntMethodCallBegin)
			if info.input.IsMatch(input) {
				return info.return0, info.return1
			}
		}
	}

	return return0, return1
}

// TODO: come up with a better name!
type Expectation struct {
	method string
	info   interface{}
}

func (m *Mock) Setup(expectation Expectation) {
	m.expectations = append([]Expectation{expectation}, m.expectations...)
}

type AddMethodCallBegin struct {
	a        kelpie.Matcher[int]
	b        kelpie.Matcher[int]
	result   int
	panicArg any
	callback func(a, b int) int
}

func Add[A1 int | kelpie.Matcher[int], A2 int | kelpie.Matcher[int]](a A1, b A2) AddMethodCallBegin {
	var a1 kelpie.Matcher[int]
	var a2 kelpie.Matcher[int]

	if matcher, ok := any(a).(kelpie.Matcher[int]); ok {
		a1 = matcher
	} else {
		a1 = kelpie.ExactMatch(any(a).(int))
	}

	if matcher, ok := any(b).(kelpie.Matcher[int]); ok {
		a2 = matcher
	} else {
		a2 = kelpie.ExactMatch(any(b).(int))
	}

	return AddMethodCallBegin{
		a: a1,
		b: a2,
	}
}

func (a AddMethodCallBegin) Return(result int) Expectation {
	a.result = result

	return Expectation{
		method: "Add",
		info:   a,
	}
}

func (a AddMethodCallBegin) Panic(arg any) Expectation {
	a.panicArg = arg

	return Expectation{
		method: "Add",
		info:   a,
	}
}

func (a AddMethodCallBegin) When(callback func(a, b int) int) Expectation {
	a.callback = callback

	return Expectation{
		method: "Add",
		info:   a,
	}
}

type ParseIntMethodCallBegin struct {
	input   kelpie.Matcher[string]
	return0 int
	return1 error
}

func ParseInt[A1 string | kelpie.Matcher[string]](input A1) ParseIntMethodCallBegin {
	var a1 kelpie.Matcher[string]

	if matcher, ok := any(input).(kelpie.Matcher[string]); ok {
		a1 = matcher
	} else {
		a1 = kelpie.ExactMatch(any(input).(string))
	}

	return ParseIntMethodCallBegin{
		input: a1,
	}
}

func (a ParseIntMethodCallBegin) Return(return0 int, return1 error) Expectation {
	a.return0 = return0
	a.return1 = return1

	return Expectation{
		method: "ParseInt",
		info:   a,
	}
}
