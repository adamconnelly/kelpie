package main

import (
	"errors"
	"fmt"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocks/emailservice"
	"github.com/adamconnelly/kelpie/mocks/maths"
)

type Maths interface {
	Add(a, b int) int
	ParseInt(input string) (int, error)
}

type EmailService interface {
	Send(sender, recipient, body string) (cost float64, err error)
}

func main() {
	emailServiceMock := emailservice.Mock{}
	emailServiceMock.Setup(emailservice.Send(kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string]()).Return(0, errors.New("unknown sender!")))
	emailServiceMock.Setup(emailservice.Send("panic@thedisco.com", kelpie.Any[string](), "testing").Panic("Oh shit!"))
	emailServiceMock.Setup(emailservice.Send("adam@email.com", "someone@receiver.com", "Hello world!").Return(100.54, nil))

	sendResult0, sendResultErr0 := emailServiceMock.Send("adam@email.com", "someone@receiver.com", "Hello world!")
	fmt.Printf("Send result 0: %f, err: %v\n", sendResult0, sendResultErr0)

	sendResult1, sendResultErr1 := emailServiceMock.Send("a", "b", "c")
	fmt.Printf("Send result 1: %f, err: %v\n", sendResult1, sendResultErr1)

	emailServiceMock.Send("panic@thedisco.com", "abc", "testing")

	mock := maths.Mock{}
	var m Maths = &mock

	// Basic expectation with exact matching
	mock.Setup(maths.Add(3, 4).Return(7))
	result1 := m.Add(3, 4)
	fmt.Printf("result1: %d\n", result1)

	// No predefined setup, so we just return the default
	result2 := m.Add(2, 6)
	fmt.Printf("result2: %d\n", result2)

	// Multiple return args
	mock.Setup(maths.ParseInt("5").Return(5, nil))
	r, err := m.ParseInt("5")
	fmt.Printf("ParseInt: %d, %v\n", r, err)

	// Mocking errors
	mock.Setup(maths.ParseInt("zzz").Return(0, errors.New("could not parse int! :(")))
	r, err = m.ParseInt("zzz")
	fmt.Printf("ParseInt: %d, %v\n", r, err)

	// No predefined setup, so we just return the default
	r, err = m.ParseInt("xxx")
	fmt.Printf("ParseInt: %d, %v\n", r, err)

	// Using a callback to intercept the call
	mock.Setup(maths.Add(10, 20).When(func(a, b int) int {
		return a + b
	}))

	result3 := m.Add(10, 20)
	fmt.Printf("Callback add: %d\n", result3)

	// Using Any matching
	mock.Setup(maths.Add(kelpie.Any[int](), kelpie.Any[int]()).Return(42))

	answer := m.Add(555, 666)
	fmt.Printf("Any result: %d\n", answer)

	// Resetting the mock and removing any expectation
	mock.Reset()

	// Using a callback to match an argument
	mock.Setup(maths.ParseInt(kelpie.Match[string](func(s string) bool {
		return s == "123"
	})).Return(123, nil))

	callbackArgResult1, _ := m.ParseInt("123")
	fmt.Printf("Callback result: %d\n", callbackArgResult1)

	callbackArgResult2, _ := m.ParseInt("321")
	fmt.Printf("Callback result 2\n: %d", callbackArgResult2)

	mock.Setup(maths.Add(kelpie.Any[int](), kelpie.Any[int]()).Return(0))

	// Mocking a PANIC!!!
	mock.Setup(maths.Add(7, 9).Panic("OMG!!!"))

	m.Add(7, 9)
}
