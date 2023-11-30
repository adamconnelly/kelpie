package mocking_test

import (
	"errors"
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocking"
	"github.com/adamconnelly/kelpie/nullable"
	"github.com/stretchr/testify/suite"
)

type MockTests struct {
	suite.Suite
}

type position struct {
	x, y, z float64
}

type fakeExpectationCreator struct {
	expectation *mocking.Expectation
}

func (c *fakeExpectationCreator) CreateExpectation() *mocking.Expectation {
	return c.expectation
}

func wrapExpectation(expectation *mocking.Expectation) *fakeExpectationCreator {
	return &fakeExpectationCreator{
		expectation: expectation,
	}
}

type fakeMethodMatcherCreator struct {
	methodMatcher *mocking.MethodMatcher
}

func (c *fakeMethodMatcherCreator) CreateMethodMatcher() *mocking.MethodMatcher {
	return c.methodMatcher
}

func wrapMethodMatcher(matcher *mocking.MethodMatcher) *fakeMethodMatcherCreator {
	return &fakeMethodMatcherCreator{
		methodMatcher: matcher,
	}
}

func (t *MockTests) TestCall_RecordsMethodCalls() {
	// Arrange
	mock := mocking.Mock{}

	// Act
	mock.Call("Launch")
	mock.Call("IncreaseVelocity", 20)
	mock.Call("SetTarget", position{x: 1, y: 2, z: 3})
	mock.Call("SetWaypoints", position{x: 1, y: 2, z: 3}, position{x: 3, y: 2, z: 1})

	// Assert
	t.Len(mock.MethodCalls, 4)
	t.Equal(*mock.MethodCalls[0], mocking.MethodCall{MethodName: "Launch", Args: nil})
	t.Equal(*mock.MethodCalls[1], mocking.MethodCall{MethodName: "IncreaseVelocity", Args: []any{20}})
	t.Equal(*mock.MethodCalls[2], mocking.MethodCall{MethodName: "SetTarget", Args: []any{position{x: 1, y: 2, z: 3}}})
	t.Equal(*mock.MethodCalls[3], mocking.MethodCall{MethodName: "SetWaypoints", Args: []any{position{x: 1, y: 2, z: 3}, position{x: 3, y: 2, z: 1}}})
}

func (t *MockTests) TestCall_ReturnsMatchingExpectation() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](-20)},
		},
		Returns: []any{errors.New("nope")},
	}))

	// Act
	expectation := mock.Call("IncreaseVelocity", -20)

	// Assert
	t.NotNil(expectation)
	t.EqualError(expectation.Returns[0].(error), "nope")
}

func (t *MockTests) TestCall_ReturnsMostRecentlySetupExpectation() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](20)},
		},
		Returns: []any{errors.New("nope")},
	}))
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](20)},
		},
		Returns: []any{nil},
	}))

	// Act
	expectation := mock.Call("IncreaseVelocity", 20)

	// Assert
	t.NotNil(expectation)
	t.Len(expectation.Returns, 1)
	t.Nil(expectation.Returns[0])
}

func (t *MockTests) TestCall_ReturnsFirstMatchingCall() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](-20)},
		},
		Returns: []any{errors.New("nope")},
	}))
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](12345)},
		},
		Returns: []any{nil},
	}))
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "ReduceVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](12345)},
		},
		Returns: []any{errors.New("nope")},
	}))
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](-40)},
		},
		Returns: []any{errors.New("nope")},
	}))

	// Act
	expectation := mock.Call("IncreaseVelocity", 12345)

	// Assert
	t.NotNil(expectation)
	t.Len(expectation.Returns, 1)
	t.Nil(expectation.Returns[0])
}

func (t *MockTests) TestCall_PanicsIfParameterCountDoesNotMatch() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](123), kelpie.ExactMatch[int](321)},
		},
		Returns: []any{errors.New("nope")},
	}))

	// Act / Assert
	t.PanicsWithValue(
		"Argument mismatch in call to 'IncreaseVelocity'.\n    Expected: 2\n    Actual: 1\nThis is a bug in Kelpie - please report it!",
		func() {
			mock.Call("IncreaseVelocity", 12345)
		})
}

func (t *MockTests) TestCall_MatchesExpectationUnlimitedTimesByDefault() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](20)},
		},
		Returns: []any{errors.New("nope")},
	}))

	// Act
	expectation1 := mock.Call("IncreaseVelocity", 20)
	expectation2 := mock.Call("IncreaseVelocity", 20)
	expectation3 := mock.Call("IncreaseVelocity", 20)

	// Assert
	t.NotNil(expectation1)
	t.NotNil(expectation2)
	t.NotNil(expectation3)
}

func (t *MockTests) TestCall_MatchesExpectationSpecifiedNumberOfTimes() {
	// Arrange
	mock := mocking.Mock{}
	mock.Setup(wrapExpectation(&mocking.Expectation{
		MethodMatcher: &mocking.MethodMatcher{
			MethodName:       "IncreaseVelocity",
			ArgumentMatchers: []mocking.ArgumentMatcher{kelpie.ExactMatch[int](20)},
			Times:            nullable.OfValue(3),
		},
		Returns: []any{errors.New("nope")},
	}))

	// Act
	expectation1 := mock.Call("IncreaseVelocity", 20)
	expectation2 := mock.Call("IncreaseVelocity", 20)
	expectation3 := mock.Call("IncreaseVelocity", 20)
	expectation4 := mock.Call("IncreaseVelocity", 20)

	// Assert
	t.NotNil(expectation1)
	t.NotNil(expectation2)
	t.NotNil(expectation3)
	t.Nil(expectation4)
}

func (t *MockTests) TestCalled_ReturnsFalseIfNoMethodsHaveBeenCalled() {
	// Arrange
	mock := mocking.Mock{}

	// Act
	called := mock.Called(wrapMethodMatcher(&mocking.MethodMatcher{}))

	// Assert
	t.False(called)
}

func (t *MockTests) TestCalled_ReturnsTrueIfMatchingCallIsFound() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20)

	// Act
	called := mock.Called(
		wrapMethodMatcher(
			&mocking.MethodMatcher{
				MethodName: "IncreaseVelocity",
				ArgumentMatchers: []mocking.ArgumentMatcher{
					kelpie.ExactMatch[int](20),
				}}))

	// Assert
	t.True(called)
}

func (t *MockTests) TestCalled_ReturnsFalseIfNoMatchingCallIsFound() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20)
	mock.Call("TrainDragon")
	mock.Call("SetTarget", position{x: 1, y: 2, z: 3})

	// Act
	called := mock.Called(
		wrapMethodMatcher(
			&mocking.MethodMatcher{
				MethodName: "Explode",
				ArgumentMatchers: []mocking.ArgumentMatcher{
					kelpie.ExactMatch[int](20),
				}}))

	// Assert
	t.False(called)
}

func (t *MockTests) TestCalled_ReturnsFalseIfNotAllParametersMatch() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20, 30)

	// Act
	called := mock.Called(
		wrapMethodMatcher(
			&mocking.MethodMatcher{
				MethodName: "IncreaseVelocity",
				ArgumentMatchers: []mocking.ArgumentMatcher{
					kelpie.ExactMatch[int](1),
					kelpie.ExactMatch[int](30),
				}}))

	// Assert
	t.False(called)
}

func (t *MockTests) TestCalled_ReturnsFalseIfMethodNotCalledEnoughTimes() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20, 30)
	mock.Call("IncreaseVelocity", 20, 30)
	mock.Call("IncreaseVelocity", 20, 30)

	// Act
	called := mock.Called(
		wrapMethodMatcher(
			&mocking.MethodMatcher{
				MethodName: "IncreaseVelocity",
				ArgumentMatchers: []mocking.ArgumentMatcher{
					kelpie.Any[int](),
					kelpie.Any[int](),
				},
				Times: nullable.OfValue(4),
			}))

	// Assert
	t.False(called)
}

func (t *MockTests) TestCalled_ReturnsTrueIfMethodCalledEnoughTimes() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20, 30)
	mock.Call("IncreaseVelocity", 20, 30)
	mock.Call("IncreaseVelocity", 20, 30)

	// Act
	called := mock.Called(
		wrapMethodMatcher(
			&mocking.MethodMatcher{
				MethodName: "IncreaseVelocity",
				ArgumentMatchers: []mocking.ArgumentMatcher{
					kelpie.Any[int](),
					kelpie.Any[int](),
				},
				Times: nullable.OfValue(3),
			}))

	// Assert
	t.True(called)
}

func (t *MockTests) TestCalled_PanicsIfArgumentCountDoesNotMatch() {
	// Arrange
	mock := mocking.Mock{}
	mock.Call("IncreaseVelocity", 20)

	// Act
	// Act / Assert
	t.PanicsWithValue(
		"Argument mismatch in call to 'IncreaseVelocity'.\n    Expected: 3\n    Actual: 1\nThis is a bug in Kelpie - please report it!",
		func() {
			mock.Called(
				wrapMethodMatcher(
					&mocking.MethodMatcher{
						MethodName: "IncreaseVelocity",
						ArgumentMatchers: []mocking.ArgumentMatcher{
							kelpie.ExactMatch[int](20),
							kelpie.ExactMatch[int](30),
							kelpie.ExactMatch[int](40),
						}}))
		})
}

func TestMock(t *testing.T) {
	suite.Run(t, new(MockTests))
}
