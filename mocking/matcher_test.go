package mocking_test

import (
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocking"
	"github.com/stretchr/testify/suite"
)

type MatcherTests struct {
	suite.Suite
}

type options struct {
	name *string
}

func WithName(name string) func(opts *options) {
	return func(opts *options) {
		opts.name = &name
	}
}

func (t *MatcherTests) Test_VariadicMatcher_MatchesWhenParametersEmpty() {
	type testCase struct {
		matchers []mocking.ArgumentMatcher
		inputs   any
		isMatch  bool
	}

	testCases := map[string]testCase{
		"Empty": {
			matchers: []mocking.ArgumentMatcher{},
			inputs:   []any{},
			isMatch:  true,
		},
		"Matchers empty but arguments provided": {
			matchers: []mocking.ArgumentMatcher{},
			inputs:   []any{"testing", 1, 2, 3},
			isMatch:  false,
		},
		"Matchers provided by arguments empty": {
			matchers: []mocking.ArgumentMatcher{kelpie.ExactMatch("testing")},
			inputs:   []any{},
			isMatch:  false,
		},
		"Arguments match": {
			matchers: []mocking.ArgumentMatcher{kelpie.ExactMatch("testing"), kelpie.ExactMatch(1), kelpie.ExactMatch(2), kelpie.ExactMatch(3)},
			inputs:   []any{"testing", 1, 2, 3},
			isMatch:  true,
		},
		"Arguments do not match": {
			matchers: []mocking.ArgumentMatcher{kelpie.ExactMatch("testing"), kelpie.ExactMatch(1), kelpie.ExactMatch(2), kelpie.ExactMatch(3)},
			inputs:   []any{"testing", 3, 2, 1},
			isMatch:  false,
		},
		"Arguments are not a slice": {
			matchers: []mocking.ArgumentMatcher{kelpie.ExactMatch("testing"), kelpie.ExactMatch(1), kelpie.ExactMatch(2), kelpie.ExactMatch(3)},
			inputs:   "testing",
			isMatch:  false,
		},
		"Arguments are slice of non-any type": {
			matchers: []mocking.ArgumentMatcher{kelpie.Any[func(*options)]()},
			inputs:   []func(*options){WithName("Bob")},
			isMatch:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func() {
			// Arrange
			matcher := mocking.Variadic(tc.matchers)

			// Act
			isMatch := matcher.IsMatch(tc.inputs)

			// Assert
			t.Equal(tc.isMatch, isMatch)
		})
	}
}

func TestMatchers(t *testing.T) {
	suite.Run(t, new(MatcherTests))
}
