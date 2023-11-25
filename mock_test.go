package kelpie_test

import (
	"errors"
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/stretchr/testify/suite"
)

type MockTests struct {
	suite.Suite
}

type position struct {
	x, y, z float64
}

func (t *MockTests) TestCall_RecordsMethodCalls() {
	// Arrange
	mock := kelpie.Mock{}

	// Act
	mock.Call("Launch")
	mock.Call("IncreaseVelocity", 20)
	mock.Call("SetTarget", position{x: 1, y: 2, z: 3})
	mock.Call("SetWaypoints", position{x: 1, y: 2, z: 3}, position{x: 3, y: 2, z: 1})

	// Assert
	t.Len(mock.MethodCalls, 4)
	t.Equal(*mock.MethodCalls[0], kelpie.MethodCall{MethodName: "Launch", Args: nil})
	t.Equal(*mock.MethodCalls[1], kelpie.MethodCall{MethodName: "IncreaseVelocity", Args: []any{20}})
	t.Equal(*mock.MethodCalls[2], kelpie.MethodCall{MethodName: "SetTarget", Args: []any{position{x: 1, y: 2, z: 3}}})
	t.Equal(*mock.MethodCalls[3], kelpie.MethodCall{MethodName: "SetWaypoints", Args: []any{position{x: 1, y: 2, z: 3}, position{x: 3, y: 2, z: 1}}})
}

func (t *MockTests) TestCall_ReturnsMatchingExpectation() {
	// Arrange
	mock := kelpie.Mock{}
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](-20)},
		Returns:          []any{errors.New("nope")},
	})

	// Act
	expectation := mock.Call("IncreaseVelocity", -20)

	// Assert
	t.NotNil(expectation)
	t.EqualError(expectation.Returns[0].(error), "nope")
}

func (t *MockTests) TestCall_ReturnsMostRecentlySetupExpectation() {
	// Arrange
	mock := kelpie.Mock{}
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](20)},
		Returns:          []any{errors.New("nope")},
	})
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](20)},
		Returns:          []any{nil},
	})

	// Act
	expectation := mock.Call("IncreaseVelocity", 20)

	// Assert
	t.NotNil(expectation)
	t.Len(expectation.Returns, 1)
	t.Nil(expectation.Returns[0])
}

func (t *MockTests) TestCall_ReturnsFirstMatchingCall() {
	// Arrange
	mock := kelpie.Mock{}
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](-20)},
		Returns:          []any{errors.New("nope")},
	})
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](12345)},
		Returns:          []any{nil},
	})
	mock.Setup(&kelpie.Expectation{
		MethodName:       "ReduceVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](12345)},
		Returns:          []any{errors.New("nope")},
	})
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](-40)},
		Returns:          []any{errors.New("nope")},
	})

	// Act
	expectation := mock.Call("IncreaseVelocity", 12345)

	// Assert
	t.NotNil(expectation)
	t.Len(expectation.Returns, 1)
	t.Nil(expectation.Returns[0])
}

func (t *MockTests) TestCall_PanicsIfParameterCountDoesNotMatch() {
	// Arrange
	mock := kelpie.Mock{}
	mock.Setup(&kelpie.Expectation{
		MethodName:       "IncreaseVelocity",
		ArgumentMatchers: []kelpie.ArgumentMatcher{kelpie.ExactMatch[int](123), kelpie.ExactMatch[int](321)},
		Returns:          []any{errors.New("nope")},
	})

	// Act / Assert
	t.PanicsWithValue(
		"Argument mismatch in call to 'IncreaseVelocity'.\n    Expected: 2\n    Actual: 1\nThis is a bug in Kelpie - please report it!",
		func() {
			mock.Call("IncreaseVelocity", 12345)
		})
}

func TestMock(t *testing.T) {
	suite.Run(t, new(MockTests))
}
