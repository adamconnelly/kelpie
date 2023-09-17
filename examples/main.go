package main

import (
	"errors"
	"fmt"
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/packages"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocks/maths"
)

type Maths interface {
	Add(a, b int) int
	ParseInt(input string) (int, error)
}

func main() {
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

	// TODO: play with a generic interface

	config := &packages.Config{
		Mode: packages.NeedSyntax,
		Dir:  "/home/adam/github.com/adamconnelly/go-better-mocks",
	}
	pkgs, _ := packages.Load(config, "adamconnelly/go-better-mocks/idea1/main")
	pkg := pkgs[0]

	for _, s := range pkg.Syntax {
		for n, o := range s.Scope.Objects {
			if o.Kind == ast.Typ {
				// check if type is exported(only need for non-local types)
				if unicode.IsUpper([]rune(n)[0]) {
					// note that reflect.ValueOf(*new(%s)) won't work with interfaces
					fmt.Printf("ProcessType(new(package_name.%s)),\n", n)
				}
			}
		}
	}

	mock.Setup(maths.Add(kelpie.Any[int](), kelpie.Any[int]()).Return(0))

	// Mocking a PANIC!!!
	mock.Setup(maths.Add(7, 9).Panic("OMG!!!"))

	m.Add(7, 9)
}
