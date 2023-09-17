package kelpie

type Matcher[T comparable] struct {
	isMatch    func(input T) bool
	exactMatch T
	isAny      bool
}

func ExactMatch[T comparable](exactMatch T) Matcher[T] {
	return Matcher[T]{exactMatch: exactMatch}
}

func (i Matcher[T]) IsMatch(test T) bool {
	return (i.isMatch != nil && i.isMatch(test)) || i.exactMatch == test || i.isAny
}

func Any[T comparable]() Matcher[T] {
	return Matcher[T]{isAny: true}
}

func Match[T comparable](isMatch func(arg T) bool) Matcher[T] {
	return Matcher[T]{isMatch: isMatch}
}
